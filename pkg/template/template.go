package template

import (
	"fmt"
	"time"
	"bytes"
	tpl "html/template"
)

type DataItem struct {
	Title    string
	Link     string
	Content  string
}

type Data struct {
	Title            string
	UnsubscribeLink  string
	Items            []*DataItem
}

type EmailTemplate struct {
	ID         int64
	ProjectID  int64
	Name 	   string
	IsActive   bool
	Subject    string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func New() *EmailTemplate {
	return &EmailTemplate{}
}

// Create the EmailTemplate record in the database
func (et *EmailTemplate) Create() error {
	if err := insertEmailTemplate(et); err != nil {
		return fmt.Errorf("Failed creating EmailTemplate: %v", err)
	}

	return nil
}

// Update the EmailTemplate record in the database
func (et *EmailTemplate) Update() error {
	if err := updateEmailTemplate(et); err != nil {
		return fmt.Errorf("Failed updating EmailTemplate: %v", err)
	}

	return nil
}

// Delete the EmailTemplate from the database
func (et *EmailTemplate) Delete() error {
	if err := deleteEmailTemplate(et.ID); err != nil {
		return fmt.Errorf("Failed deleting EmailTemplate: %v", err)
	}

	return nil
}

// Activate sets the attribute IsActive to true and updates the EmailTemplate record
func (et *EmailTemplate) Activate() error {
	et.IsActive = true
	if err := et.Update(); err != nil {
		return fmt.Errorf("Failed activating EmailTemplate: %v", err)
	}

	return nil
}

// Deactivate sets the attribute IsActive to false and updates the EmailTemplate record
func (et *EmailTemplate) Deactivate() error {
	et.IsActive = false
	if err := et.Update(); err != nil {
		return fmt.Errorf("Failed activating EmailTemplate: %v", err)
	}

	return nil
}

// RenderContent receives data to build the email content
func (et *EmailTemplate) RenderContent(data *Data) (string, error) {
	return render(et.Content, data)
}

// RenderSubject receives data to build the email subject
func (et *EmailTemplate) RenderSubject(data *Data) (string, error) {
	return render(et.Subject, data)
}

// render receives a template string and any object that matches
// variables defined in this template. Then, it builds the template using
// the data and returns the resulting string
func render(t string, data interface{}) (string, error) {
	renderer, err := tpl.New("renderer").Parse(t)
	if err != nil {
		return "", err
	}
	
	var content bytes.Buffer
	if err := renderer.Execute(&content, data); err != nil {
		return "", err
	}

	return content.String(), nil
}
