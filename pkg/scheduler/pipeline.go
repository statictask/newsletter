package scheduler

import (
	"time"

	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/pkg/project"
	"go.uber.org/zap"
)

type PipelineScheduler struct{}

func NewPipelineScheduler() *PipelineScheduler {
	return &PipelineScheduler{}
}

// Start creates a go routine to reconcile project's pipelines
func (s *PipelineScheduler) Start() (chan Signal, error) {
	ch := make(chan Signal)

	go s.startPipelineReconcileLoop(ch)

	return ch, nil
}

// startPipelineReconcileLoop checks whether there's a condition
// in which the scheduler needs to create a new pipeline for projects enabled
func (s *PipelineScheduler) startPipelineReconcileLoop(stop chan Signal) {
	log.L.Info("project pipeline reconcile loop started")
	for {
		time.Sleep(10 * time.Second)

		select {
		case <-stop:
			log.L.Info("project pipeline reconcile loop stopped")
			return

		default:
			projects := project.NewProjects()
			enabledProjects, err := projects.AllEnabled()
			if err != nil {
				log.L.Error("project pipeline reconcile loop failed to get enabled projects", zap.Error(err))
				continue
			}

			for _, p := range enabledProjects {
				log.L.Info("reconciling project pipelines", zap.Int64("project_id", p.ID))

				lastPipeline, err := p.Pipelines().Last()
				if err != nil {
					log.L.Error("failed getting project pipeline", zap.Error(err))
					continue
				}

				if lastPipeline == nil {
					log.L.Info("creating the first project pipeline", zap.Int64("project_id", p.ID))

					np, err := p.Pipelines().Create()
					if err != nil {
						log.L.Error("failed creating new pipeline", zap.Int64("project_id", p.ID), zap.Error(err))
					}

					log.L.Info("new pipeline created", zap.Int64("pipeline_id", np.ID), zap.Int64("project_id", p.ID), zap.Error(err))
					continue
				}

				is_finished, err := lastPipeline.IsFinished()
				if err != nil {
					log.L.Error("failed checking pipeline status", zap.Error(err))
					continue
				}

				// If the pipeline is not finished we don't do anything
				if !is_finished {
					log.L.Info("skipping pipeline creation", zap.Int64("project_id", p.ID))
					continue
				}

				log.L.Info("creating a new project pipeline", zap.Int64("project_id", p.ID))
				np, err := p.Pipelines().Create()
				if err != nil {
					log.L.Error("failed creating new pipeline", zap.Int64("project_id", p.ID), zap.Error(err))
				}

				log.L.Info("new pipeline created", zap.Int64("pipeline_id", np.ID), zap.Int64("project_id", p.ID), zap.Error(err))
			}
		}
	}
}
