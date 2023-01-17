package publisher

type EmailAddress struct {
	Name string
	Address string
}

type Email struct {
	From *EmailAddress
	To *EmailAddress
	Subject string
	Content string

}

func NewEmail(from, to *EmailAddress, subject, content string) *Email {
	return &Email{from, to, subject, content}
}

func NewEmailAddress(name, address string) *EmailAddress {
	return &EmailAddress{name, address}
}
