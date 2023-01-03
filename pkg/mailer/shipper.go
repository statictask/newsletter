package mailer

import (
	"net/smtp"
)

type Shipper struct{
	Server string
	User string
	Password string
}

func NewShipper(server, user, password string) *Shipper {
	return &Shipper{server, user, password}
}

func (s *Shipper) Send(mail *Mail) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", s.User, s.Password, s.Server)
	from := mail.From
	to := []string{mail.To}
	msg := []byte("To: " + mail.To + "\r\n" +
		"Subject: " + mail.Subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + mail.Content,
	)

	err := smtp.SendMail(s.Server + ":25", auth, from, to, msg)
	if err != nil {
		return err
	}

	return nil
}
