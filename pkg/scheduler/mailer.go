package scheduler

import "github.com/statictask/newsletter/pkg/publisher"

type PublisherJobScheduler struct{}

func NewPublisherJobScheduler() *PublisherJobScheduler {
	return &PublisherJobScheduler{}
}

// Start creates a go routine to reconcile pipeline's tasks
func (s *PublisherJobScheduler) Start() {
	job := publisher.New()
	go job.Run()
}
