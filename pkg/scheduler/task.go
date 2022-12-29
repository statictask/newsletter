package scheduler

import (
	"time"

	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/pkg/pipeline"
	"github.com/statictask/newsletter/pkg/project"
	"go.uber.org/zap"
)

type TaskScheduler struct{}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{}
}

// Start creates a go routine to reconcile pipeline's tasks
func (s *TaskScheduler) Start() (chan Signal, error) {
	ch := make(chan Signal)

	go s.startTaskReconcileLoop(ch)

	return ch, nil
}

// startTaskReconcileLoop checks whether there's a condition
// in which the scheduler needs to create a new task for existing pipelines
func (s *TaskScheduler) startTaskReconcileLoop(stop chan Signal) {
	log.L.Info("task reconcile loop started")
	for {
		time.Sleep(10 * time.Second)

		select {
		case <-stop:
			log.L.Info("task reconcile loop stopped")
			return

		default:
			projects := project.NewProjects()
			enabledProjects, err := projects.AllEnabled()
			if err != nil {
				log.L.Error("task reconcile loop failed to get enabled projects", zap.Error(err))
				continue
			}

			for _, p := range enabledProjects {
				log.L.Info("reconciling tasks", zap.Int64("project_id", p.ID))

				lastPipeline, err := p.Pipelines().Last()
				if err != nil {
					log.L.Error("failed getting project pipeline", zap.Error(err))
					continue
				}

				if lastPipeline == nil {
					log.L.Info("project does not have available pipelines", zap.Int64("project_id", p.ID))
					continue
				}

				for _, taskType := range pipeline.TaskTypes {
					task, err := lastPipeline.Tasks().GetByType(pipeline.Scrape)
					if err != nil {
						log.L.Error(
							"failed loading pipeline task",
							zap.Int64("project_id", p.ID),
							zap.Int64("pipeline_id", lastPipeline.ID),
							zap.Reflect("task_type", taskType),
							zap.Error(err),
						)
						continue
					}

					if task == nil {
						log.L.Info(
							"creating a new task",
							zap.Int64("project_id", p.ID),
							zap.Int64("pipeline_id", lastPipeline.ID),
						)

						task, err := lastPipeline.Tasks().Create(taskType)
						if err != nil {
							log.L.Error(
								"failed creating new task",
								zap.Int64("project_id", p.ID),
								zap.Int64("pipeline_id", lastPipeline.ID),
								zap.Reflect("task_type", taskType),
								zap.Error(err),
							)
						}

						log.L.Info(
							"new task created",
							zap.Int64("project_id", p.ID),
							zap.Int64("pipeline_id", lastPipeline.ID),
							zap.Reflect("task_type", task.Type),
							zap.Error(err),
						)
					}
				}
			}
		}
	}
}
