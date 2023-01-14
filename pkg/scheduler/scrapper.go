package scheduler

import (
	"time"
	"github.com/statictask/newsletter/pkg/scrapper"
)

type ScrapperJobScheduler struct{}

func NewScrapperJobScheduler() *ScrapperJobScheduler {
	return &ScrapperJobScheduler{}
}

// Start creates a go routine to reconcile pipeline's tasks
func (s *ScrapperJobScheduler) Start(minScrapeInterval time.Duration, allowPreviousPublications bool) {
	job := scrapper.New(minScrapeInterval, allowPreviousPublications)
	go job.Run()
}
