package mailer

import (
	"time"
	"fmt"
	"context"

	"go.uber.org/zap"

	"github.com/statictask/newsletter/pkg/task"
	"github.com/statictask/newsletter/pkg/post"
	"github.com/statictask/newsletter/pkg/project"
	"github.com/statictask/newsletter/pkg/subscription"
	"github.com/statictask/newsletter/internal/log"
)

type EmailSender interface {
	Send (ctx context.Context, tm *TargetEmail) error
}

type Mailer struct {
	sender EmailSender
}

func New(name, address, key string) *Mailer {
	// define mail shipper
	addr := NewEmailAddress(name, address)
	shipper := NewShipper(addr, key)

	return &Mailer{shipper}
}

// Run executes an infinite loop that keeps checking if there are
// new posts to be sent
func (m *Mailer) Run() {
	for {
		time.Sleep(10 * time.Second)

		if err := m.processWaitingTasks(); err != nil {
			log.L.Error("failed processing publish waiting task", zap.Error(err))
		}

		if err := m.processReadyTasks(); err != nil {
			log.L.Error("failed processing publish ready tasks", zap.Error(err))
		}
	}
}

func (m *Mailer) processWaitingTasks() error {
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
			_log.Error("failed loading pipeline related task", zap.Error(err))
			continue
		}

		if scrapeTask.Status != task.Finished {
			_log.Debug("scrape task is not finished", zap.Error(err))
			continue
		}
		
		t.Status = task.Ready
		if err := t.Update(); err != nil {
			_log.Error("failed updating task", zap.Error(err))
			continue
		}

		_log.Info("publish task is ready to be processed")
	}

	return nil
}

func (m *Mailer) processReadyTasks() error {
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
		relatedPost, err := pipelinePosts.Last()
		if err != nil {
			_log.Error("pipeline post not found", zap.Error(err))

			t.Status = task.Failed
			if err := t.Update(); err != nil {
				_log.Error("failed to mark task as failed", zap.Error(err))
				continue
			}
		}

		taskProject, err := project.NewProjects().GetByTaskID(t.ID)
		if err != nil {
			_log.Error("failed getting the project of the task", zap.Error(err))
			continue
		}

		_log = _log.With(zap.Int64("project_id", taskProject.ID))

		subscriptions, err := taskProject.Subscriptions().All()
		if err != nil {
			_log.Error("failed getting task project subscriptions", zap.Error(err))
			continue
		}

		deliveryCount := 0
		deliveryTotal := len(subscriptions)

		for _, s := range subscriptions {
			__log := _log.With(zap.Int64("subscription_id", s.ID))

			if err := m.sendEmail(s, relatedPost); err != nil {
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

func (m *Mailer) sendEmail(s *subscription.Subscription, p *post.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// create email
	mail := NewTargetEmail(s, p)

	return m.sender.Send(ctx, mail)
}
