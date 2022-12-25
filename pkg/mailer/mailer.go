package mailer

import (
	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

type Mailer struct{}

func New() *Mailer {
	return &Mailer{}
}

func (m *Mailer) Run() chan string {
	message := make(chan string)

	go func() {
		for m := range message {
			log.L.Info("new mail arrieved", zap.String("message", m))
		}
	}()

	return message
}
