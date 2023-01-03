package scrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/pkg/post"
	"github.com/statictask/newsletter/pkg/project"
	"github.com/statictask/newsletter/pkg/task"
	"go.uber.org/zap"
)

type Scrapper struct{}

func New() *Scrapper {
	return &Scrapper{}
}

// Run checks if there are new items published in the feed since the last
// finished pipeline. If so, it'll get these items and create a post with
// the new content in the database
func (s *Scrapper) Run() {
	for {
		time.Sleep(10 * time.Second)

		if err := s.processWaitingTasks(); err != nil {
			log.L.Error("failed processing scrape waiting tasks", zap.Error(err))
		}

		if err := s.processReadyTasks(); err != nil {
			log.L.Error("failed processing scrape ready tasks", zap.Error(err))
		}
	}
}

func (s *Scrapper) processWaitingTasks() error {
	ctl := task.NewTasks()
	tasks, err := ctl.Filter(task.Scrape, task.Waiting)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		t.Status = task.Ready
		if err := t.Update(); err != nil {
			log.L.Error("failed getting task ready", zap.Error(err), zap.Int64("task_id", t.ID))
			continue
		}

		log.L.Info("scrape task is ready to be processed", zap.Int64("task_id", t.ID))
	}

	return nil
}

func (s *Scrapper) processReadyTasks() error {
	ctl := task.NewTasks()
	tasks, err := ctl.Filter(task.Scrape, task.Ready)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		taskProject, err := project.NewProjects().GetByTaskID(t.ID)
		if err != nil {
			log.L.Error("failed getting the project of the task", zap.Error(err), zap.Int64("task_id", t.ID))
			continue
		}

		feedURL := taskProject.FeedURL
		feedReader := NewFeedReader(feedURL)

		lastPost, err := taskProject.Posts().Last()
		if err != nil {
			log.L.Error("failed getting the last post of project the project", zap.Error(err), zap.Int64("task_id", t.ID))
			continue
		}

		// if there are no posts, we're going to use the UpdatedAt field
		// of the project sync further posts
		isFirst := true
		lastPubDate := taskProject.UpdatedAt
		if lastPost != nil {
			isFirst = false
			lastPubDate = lastPost.UpdatedAt
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		items, err := feedReader.ReadFrom(ctx, lastPubDate, isFirst)
		cancel()

		_log := log.L.With(zap.Int64("task_id", t.ID), zap.String("feed_url", feedURL))

		if err != nil {
			_log.Info("failed reading feed", zap.Error(err))
			continue
		}

		if len(items) == 0 {
			_log.Info("nothing new on this feed")
			continue
		}

		ifaces := make([]post.ContentItem, len(items))
		for i, s := range items {
			ifaces[i] = post.ContentItem(s)
		}

		log.L.Info("found new items")

		contentBuilder := post.NewContent(
			fmt.Sprintf("The newsletter of %s", taskProject.Name),
			ifaces,
		)

		content, err := contentBuilder.Build()
		if err != nil {
			_log.Error("failed building feed content", zap.Error(err))
			continue
		}

		// TODO: check if there's a post created before adding another one
		// and raise an error in case of multiple posts
		newPost := post.New()
		newPost.PipelineID = t.PipelineID
		newPost.Content = content

		if err := newPost.Create(); err != nil {
			_log.Error("failed creating feed post", zap.Error(err))
			continue
		}

		_log.Info("successfully created feed post", zap.Int64("post_id", newPost.ID))

		t.Status = task.Finished
		if err := t.Update(); err != nil {
			_log.Error("failed building feed content", zap.Error(err))
			continue
		}

		_log.Info("finished scrape task", zap.Int64("post_id", newPost.ID))
	}

	return nil
}
