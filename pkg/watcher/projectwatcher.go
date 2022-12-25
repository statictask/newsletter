package watcher

import (
	"time"

	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/pkg/project"
	"go.uber.org/zap"
)

type WatcherSignal int64

const Stop WatcherSignal = 0

type ProjectWatcherRunner struct {
	Project *project.Project
	Watcher *FeedWatcher
	Stop    chan WatcherSignal
}

type ProjectWatcherPool struct {
	ProjectWatcherRunners map[int64]*ProjectWatcherRunner
}

func New() *ProjectWatcherPool {
	return &ProjectWatcherPool{map[int64]*ProjectWatcherRunner{}}
}

func (pwp *ProjectWatcherPool) startWatcher(target *project.Project, message chan string) {
	if _, ok := pwp.ProjectWatcherRunners[target.ID]; !ok {
		log.L.Debug("starting project feed watcher", zap.Int64("projectID", target.ID), zap.String("url", target.FeedURL))
		watcher := NewFeedWatcher(target.FeedURL)
		stop := watcher.Run(message)
		runner := &ProjectWatcherRunner{
			Project: target,
			Watcher: watcher,
			Stop:    stop,
		}

		pwp.ProjectWatcherRunners[target.ID] = runner
		log.L.Info("project feed watcher started", zap.Int64("projectID", target.ID), zap.String("url", target.FeedURL))
		return
	}

	log.L.Debug("watcher is already running", zap.Int64("projectID", target.ID), zap.String("url", target.FeedURL))
}

func (pwp *ProjectWatcherPool) stopWatcher(id int64) {
	if _, ok := pwp.ProjectWatcherRunners[id]; !ok {
		log.L.Debug("project feed watcher is not running", zap.Int64("projectID", id))
		return
	}

	log.L.Debug("stopping project feed watcher", zap.Int64("projectID", id))

	pwp.ProjectWatcherRunners[id].Stop <- Stop
	delete(pwp.ProjectWatcherRunners, id)

	log.L.Info("project feed watcher stopped", zap.Int64("projectID", id))
}

func (pwp *ProjectWatcherPool) Run(message chan string) chan WatcherSignal {
	stop := make(chan WatcherSignal)

	go func() {
		projects := project.NewProjects()

		for {
			select {
			case <-stop:
				log.L.Info("stopping all project feed watchers")
				for id, _ := range pwp.ProjectWatcherRunners {
					pwp.stopWatcher(id)
				}

				return
			default:
				log.L.Info("refreshing project feed watchers")

				allProjects, err := projects.All()
				if err != nil {
					log.L.Error("unable to load projects", zap.Error(err))
				}

				currentProjectIDs := map[int64]struct{}{}
				for _, p := range allProjects {
					pwp.startWatcher(p, message)
					currentProjectIDs[p.ID] = struct{}{}
				}

				for id, _ := range pwp.ProjectWatcherRunners {
					if _, ok := currentProjectIDs[id]; !ok {
						pwp.stopWatcher(id)
					}
				}

				time.Sleep(60 * time.Second)
			}
		}
	}()

	return stop
}
