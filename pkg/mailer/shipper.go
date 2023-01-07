package mailer

import (
	"fmt"
	"context"

	"go.uber.org/zap"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/rest"

	"github.com/statictask/newsletter/internal/log"

)

type EmailServiceClient interface {
	SendWithContext(ctx context.Context, email *mail.SGMailV3) (*rest.Response, error)
}

// Shipper uses SendGrid API to send a TargetEmail
type Shipper struct{
	Sender *EmailAddress
	Client EmailServiceClient
}

func NewShipper(sender *EmailAddress, key string) *Shipper {
	fmt.Println(key)
	client := sendgrid.NewSendClient(key)
	return &Shipper{sender, client}
}

// Send calls SendGrid API to send an email
func (s *Shipper) Send(ctx context.Context, tm *TargetEmail) error {
	from := mail.NewEmail(s.Sender.Name, s.Sender.Address)
	to := mail.NewEmail(tm.To.GetName(), tm.To.GetAddress())

	_log := log.L.With(zap.String("target", tm.To.GetAddress()))

	subject := tm.Email.GetSubject()
	plainContent, err := tm.Email.GetPlainTextContent()
	if err != nil {
		return err
	}

	htmlContent, err := tm.Email.GetHTMLContent()
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject, to, plainContent, htmlContent)
	response, err := s.Client.SendWithContext(ctx, message)
	if err != nil {
		_log.Info("failed sending email", zap.Error(err))
		return err
	} else {
		_log.Info("email successfuly sent", zap.Int("status_code", response.StatusCode))
	}

	return nil
}
