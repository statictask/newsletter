package mailer

type Email interface {
	GetSubject() string
	GetPlainTextContent() (string, error)
	GetHTMLContent() (string, error)
}

type Target interface {
	GetName() string
	GetAddress() string
}

type TargetEmail struct {
	To Target
	Email Email
}

func NewTargetEmail(to Target, email Email) *TargetEmail {
	return &TargetEmail{to, email}
}
