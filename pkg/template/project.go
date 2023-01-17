package template

import (
	"fmt"
)


type ProjectEmailTemplates struct {
	projectID int64
}

// NewProjectEmailTemplates returns a EmailTemplates controller
func NewProjectEmailTemplates(projectID int64) *ProjectEmailTemplates {
	return &ProjectEmailTemplates{projectID}
}

// All returns all project's EmailTemplates
func (pt *ProjectEmailTemplates) All() ([]*EmailTemplate, error) {
	return getEmailTemplatesByProjectID(pt.projectID)
}

// Get returns the respective EmailTemplate by ID
func (pt *ProjectEmailTemplates) Get(id int64) (*EmailTemplate, error) {
	return getEmailTemplateByID(id)
}

// GetActive returns this project's active EmailTemplate
func (pt *ProjectEmailTemplates) GetActive() (*EmailTemplate, error) {
	emailTemplate, err := getActiveEmailTemplateByProjectID(pt.projectID)
	if err != nil {
		return nil, err
	}

	if emailTemplate == nil {
		return nil, fmt.Errorf("EmailTemplate not found.")
	}

	return emailTemplate, nil
}

// Add creates a new entry in the project's email_templates 
func (pt *ProjectEmailTemplates) Add(et *EmailTemplate) error {
	// make sure the EmailTemplate has the corred ProjectID before creating
	et.ProjectID = pt.projectID

	if err := et.Create(); err != nil {
		return err
	}

	return nil
}

// CreateDefault creates a new default EmailTemplate for the project
func (pt *ProjectEmailTemplates) CreateDefault() (*EmailTemplate, error) {
	content := `<html>
		       <head>
		         <title>{{ .Title }}</title>
		       </head>
		       <body>
		         <h1>
			   {{ .Title }}
		         </h1>
		         <br>
			 {{ range .Items }}
			 <hr>
			 <a href="{{ .Link }}">
			   <h3>
			     {{ .Title }}
			   </h3>
			 </a>
			 <br>
			 <p>
			   {{ .Content }}
			 </p>
			 {{ end }}
			 <br>
			 <br>
			 <br>
			 <p>
			   This newsletter is powered by <a href="https://statictask.io">statictask.io</a>.
			 </p>
			 <p>
			   Don't want to receive this email anymore? <a href="{{ .UnsubscribeLink }}">Unsubscribe</a>.
			 </p>
			 <br>
		       </body>
		     </html>
	`

	et := New()
	et.Name = "Default"
	et.Subject = "[Newsletter] {{ .Title }}"
	et.Content = content

	if err := pt.Add(et); err != nil {
		return nil, err
	}

	return et, nil
}
