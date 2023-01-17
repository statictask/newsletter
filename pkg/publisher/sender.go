package publisher

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/rest"

	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/internal/config"

)

type EmailServiceClient interface {
	SendWithContext(ctx context.Context, email *mail.SGMailV3) (*rest.Response, error)
}

// Sender uses SendGrid API to send a TargetEmail
type Sender struct{
	Client EmailServiceClient
}

func NewSender() *Sender {
	client := sendgrid.NewSendClient(config.C.SendGridAPIKey)
	return &Sender{client}
}

// Send calls SendGrid API to send an email
func (s *Sender) Send(ctx context.Context, e *Email) error {
	_log := log.L.With(zap.String("target", e.To.Address))

	from := mail.NewEmail(e.From.Name, e.From.Address)
	to := mail.NewEmail(e.To.Name, e.To.Address)
	subject := e.Subject
	htmlContent := e.Content


	fmt.Println("------------------------")
	fmt.Println(e.Subject)
	fmt.Println("------------------------")
	fmt.Println(e.Content)
	fmt.Println("------------------------")

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	response, err := s.Client.SendWithContext(ctx, message)
	if err != nil {
		_log.Info("Failed sending email.", zap.Error(err))
		return err
	} else {
		_log.Info("Sendgrid email successfuly sent.", zap.Int("status_code", response.StatusCode))
	}

	return nil
}
