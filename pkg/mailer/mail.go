package mailer

type Mail struct {
	From string
	To string
	Subject string
	Content string
}

func NewMail(from, to, subject, content string) *Mail {
	return &Mail{from, to, subject, content}
}
