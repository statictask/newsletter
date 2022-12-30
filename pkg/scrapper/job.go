package scrapper

import (
	"time"

	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/pkg/task"
	"go.uber.org/zap"
)

type ScrapperJob struct{}

func NewScrapperJob() *ScrapperJob {
	return &ScrapperJob{}
}

// Run checks if there are new items published in the feed since the last
// finished pipeline. If so, it'll get these items and create a post with
// the new content in the database
func (s *ScrapperJob) Run() {
	for {
		time.Sleep(10 * time.Second)

		if err := s.processWaitingTasks(); err != nil {
			log.L.Error("failed processing waiting tasks", zap.Error(err))
		}

		if err := s.processReadyTasks(); err != nil {
			log.L.Error("failed processing ready tasks", zap.Error(err))
		}
	}
}

func (s *ScrapperJob) processWaitingTasks() error {
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

func (s *ScrapperJob) processReadyTasks() error {
	ctl := task.NewTasks()
	tasks, err := ctl.Filter(task.Scrape, task.Ready)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		log.L.Info("will process task", zap.Int64("task_id", t.ID))
	}

	return nil
}
