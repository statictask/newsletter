package publisher

import (
	"time"
	"fmt"
	"context"
	"net/url"

	"go.uber.org/zap"

	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/internal/config"
	"github.com/statictask/newsletter/pkg/task"
	"github.com/statictask/newsletter/pkg/post"
	"github.com/statictask/newsletter/pkg/project"
	"github.com/statictask/newsletter/pkg/subscription"
	"github.com/statictask/newsletter/pkg/template"
)

type EmailSender interface {
	Send (ctx context.Context, e *Email) error
}

type Watcher struct {
	sender EmailSender
}

func New() *Watcher {
	sender := NewSender()
	return &Watcher{sender}
}

// Run executes an infinite loop that keeps checking if there are
// new posts to be sent
func (w *Watcher) Run() {
	_log := log.L.With(zap.String("watcher", "publisher"))
	for {
		time.Sleep(10 * time.Second)

		if err := w.processWaitingTasks(); err != nil {
			_log.Error("Failed processing waiting tasks.", zap.Error(err), zap.String("stage", "waiting"))
		}

		if err := w.processReadyTasks(); err != nil {
			_log.Error("Failed processing ready tasks.", zap.Error(err), zap.String("stage", "ready"))
		}
	}
}

func (w *Watcher) processWaitingTasks() error {
	ctl := task.NewTasks()
	tasks, err := ctl.Filter(task.Publish, task.Waiting)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		_log := log.L.With(
			zap.Int64("task_id", t.ID),
			zap.String("task_type", string(t.Type)),
			zap.Int64("pipeline_id", t.PipelineID),
		)

		taskPipeline := task.NewPipelineTasks(t.PipelineID)	

		scrapeTask, err := taskPipeline.GetByType(task.Scrape)
		if err != nil {
			_log.Error("Failed loading pipeline related task.", zap.Error(err))
			continue
		}

		if scrapeTask.Status != task.Finished {
			_log.Debug("Scrape task is not finished.", zap.Error(err))
			continue
		}
		
		t.Status = task.Ready
		if err := t.Update(); err != nil {
			_log.Error("Failed updating task.", zap.Error(err))
			continue
		}

		_log.Info("Task is ready.")
	}

	return nil
}

func (w *Watcher) processReadyTasks() error {
	ctl := task.NewTasks()
	tasks, err := ctl.Filter(task.Publish, task.Ready)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		_log := log.L.With(
			zap.Int64("task_id", t.ID),
			zap.Int64("pipeline_id", t.PipelineID),
		)

		pipelinePosts := post.NewPipelinePosts(t.PipelineID)	
		lastPost, err := pipelinePosts.Last()
		if err != nil {
			_log.Error("Post not found for this Pipeline. Skipping.", zap.Error(err))

			t.Status = task.Failed
			if err := t.Update(); err != nil {
				_log.Error("Failed to mark Task as Failed. Skipping.", zap.Error(err))
				continue
			}

			continue
		}

		taskProject, err := project.NewProjects().GetByTaskID(t.ID)
		if err != nil {
			_log.Error("Failed loading the Task's Project. Skipping.", zap.Error(err))
			continue
		}

		_log = _log.With(zap.Int64("project_id", taskProject.ID))

		subscriptions, err := taskProject.Subscriptions().All()
		if err != nil {
			_log.Error("Failed loading Project's Subscriptions. Skipping.", zap.Error(err))
			continue
		}

		activeEmailTemplate, err := taskProject.EmailTemplates().GetActive()
		if err != nil {
			_log.Error("Failed loading Project's active EmailTemplate. Skipping", zap.Error(err))
			continue
		}

		deliveryCount := 0
		deliveryTotal := len(subscriptions)

		for _, s := range subscriptions {
			__log := _log.With(zap.Int64("subscription_id", s.ID))

			if err := w.sendEmail(s, lastPost, activeEmailTemplate); err != nil {
				__log.Error("failed sending email", zap.Error(err))
				continue
			}

			deliveryCount += 1
			__log.Info(fmt.Sprintf("email was sent (%d/%d)", deliveryCount, deliveryTotal))
		}

		if deliveryCount != len(subscriptions) {
			_log.Error(
				"failed to send one or more emails",
				zap.Int("failed", deliveryTotal - deliveryCount),
				zap.Int("delivered", deliveryCount),
			)

			t.Status = task.Failed
			if err := t.Update(); err != nil {
				_log.Error("failed to mark task as failed", zap.Error(err))
				continue
			}
		}

		t.Status = task.Finished
		if err := t.Update(); err != nil {
			_log.Error("failed updating task", zap.Error(err))
			continue
		}

		_log.Info("publish task is finished")
	}

	return nil
}

func (w *Watcher) sendEmail(s *subscription.Subscription, p *post.Post, et *template.EmailTemplate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Build unique links for users to unsubscribe the newsletter
	unsubscribeToken, err := s.Encrypt()
	if err != nil {
		return err
	}

	unsubscribeLink := url.URL{
		Scheme: "https",
		Host: config.C.ApplicationDomain,
		Path: "unsubscribe",
		RawQuery: fmt.Sprintf("token=%s", unsubscribeToken),
	}

	// Create post items array to be processed
	postItems, err := p.PostItems().All()
	if err != nil {
		return err
	}

	tplDataItems := []*template.DataItem{}
	for _, pi := range postItems {
		item := &template.DataItem{
			Title: pi.Title,
			Link: pi.Link,
			Content: pi.Content,
		}

		tplDataItems = append(tplDataItems, item)
	}
	
	tplData := &template.Data{
		Title: p.Title,
		UnsubscribeLink: unsubscribeLink.String(),
		Items: tplDataItems,
	}

	// Build email to be sent
	emailSubject, err := et.RenderSubject(tplData)
	if err != nil {
		return err
	}

	emailContent, err := et.RenderContent(tplData)
	if err != nil {
		return err
	}

	emailFrom := NewEmailAddress(config.C.PublisherName, config.C.PublisherEmail)
	emailTo := NewEmailAddress("Reader", s.Email)
	email := NewEmail(emailFrom, emailTo, emailSubject, emailContent)

	return w.sender.Send(ctx, email)
}
