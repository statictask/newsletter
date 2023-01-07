package mailer

type EmailAddress struct {
	Name string
	Address string
}

func NewEmailAddress(name, address string) *EmailAddress {
	return &EmailAddress{name, address}
}
