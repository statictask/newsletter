package scrapper

import (
	"context"
	"time"

	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/pkg/post"
	"github.com/statictask/newsletter/pkg/postitem"
	"github.com/statictask/newsletter/pkg/project"
	"github.com/statictask/newsletter/pkg/task"
	"go.uber.org/zap"
)

type Scrapper struct {
	AllowPreviousPublications bool
}

func New(allowPreviousPublications bool) *Scrapper {
	return &Scrapper{allowPreviousPublications}
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


		lastPost, err := taskProject.Posts().Last()
		if err != nil {
			log.L.Error("failed getting the last post of project the project", zap.Error(err), zap.Int64("task_id", t.ID))
			continue
		}

		// if there are no posts, we're going to use the UpdatedAt field
		// of the project sync further posts
		feedURL := taskProject.FeedURL
		isFirst := true
		lastPubDate := taskProject.UpdatedAt
		if lastPost != nil {
			isFirst = false
			lastPubDate = lastPost.UpdatedAt
		}

		_log := log.L.With(zap.Int64("task_id", t.ID), zap.String("feed_url", feedURL))

		items, err := s.readFeed(feedURL, lastPubDate, isFirst)
		if err != nil {
			_log.Info("failed reading feed", zap.Error(err))
			continue
		}

		if len(items) == 0 {
			_log.Info("nothing new on this feed")
			continue
		}

		log.L.Info("found new items")

		// TODO: check if there's a post created before adding another one
		// and raise an error in case of multiple posts
		newPost := post.New()
		newPost.PipelineID = t.PipelineID
		newPost.Title = "Newsletter - " + taskProject.Name

		if err := newPost.Create(); err != nil {
			_log.Error("failed creating feed post", zap.Error(err))
			continue
		}

		// TODO: clean old post items for this post before addind new ones
		// this is important in case there was an error adding post items
		// in a previous cycle and now the item may be duplicated
		_log = _log.With(zap.Int64("post_id", newPost.ID))
		_log.Info("successfully created feed post")

		for _, i := range items {
			newPostItem := postitem.New()
			newPostItem.PostID = newPost.ID
			newPostItem.Title = i.Title
			newPostItem.Link = i.Link
			newPostItem.Content = i.Content

			if err := newPostItem.Create(); err != nil {
				_log.Error("failed creating new post item", zap.Error(err))
				continue
			}
		}

		t.Status = task.Finished
		if err := t.Update(); err != nil {
			_log.Error("failed building feed content", zap.Error(err))
			continue
		}

		_log.Info("finished scrape task")
	}

	return nil
}

func (s *Scrapper) readFeed(url string, lastPubDate *time.Time, isFirst bool) ([]*FeedItem, error) {
	feedReader := NewFeedReader(url)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return feedReader.ReadFrom(ctx, lastPubDate, isFirst, s.AllowPreviousPublications)
}
