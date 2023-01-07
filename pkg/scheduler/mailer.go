package scheduler

import "github.com/statictask/newsletter/pkg/mailer"

type MailerJobScheduler struct{}

func NewMailerJobScheduler() *MailerJobScheduler {
	return &MailerJobScheduler{}
}

// Start creates a go routine to reconcile pipeline's tasks
func (s *MailerJobScheduler) Start(name, address, key string) {
	job := mailer.New(name, address, key)
	go job.Run()
}
