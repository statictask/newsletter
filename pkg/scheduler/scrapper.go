package scheduler

import "github.com/statictask/newsletter/pkg/scrapper"

type ScrapperJobScheduler struct{}

func NewScrapperJobScheduler() *ScrapperJobScheduler {
	return &ScrapperJobScheduler{}
}

// Start creates a go routine to reconcile pipeline's tasks
func (s *ScrapperJobScheduler) Start() {
	job := scrapper.New()
	go job.Run()
}
