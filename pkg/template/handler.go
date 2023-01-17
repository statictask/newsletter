package template

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/database"
)

// insertEmailTemplate creates a new row in the table email_templates
// based on the given pre-created EmailTemplate object
func insertEmailTemplate(et *EmailTemplate) error {
	query := `
		INSERT INTO email_templates (
		  project_id,
		  name,
		  subject,
		  content
	        )
		VALUES (
		  $1,
		  $2,
		  $3,
		  $4
	        )
		RETURNING
		  email_template_id,
		  project_id,
		  name,
		  is_active,
		  subject,
		  content,
		  created_at,
		  updated_at
	`

	savedEmailTemplate, err := scanEmailTemplate(
		query,
		et.ProjectID,
		et.Name,
		et.Subject,
		et.Content,
	)
	if err != nil {
		return err
	}

	*et = *savedEmailTemplate

	return nil
}

// getEmailTemplateByID returns a single email_template with the given ID if it exists
func getEmailTemplateByID(id int64) (*EmailTemplate, error) {
	query := `
		SELECT
		  email_template_id,
		  project_id,
		  name,
		  is_active,
		  subject,
		  content,
		  created_at,
		  updated_at
		FROM
		  email_templates
		WHERE
		  email_template_id = $1
	`

	return scanEmailTemplate(query, id)
}

// getEmailTemplatesByProjectID returns all email_templates in the database based
// on a given expression
func getEmailTemplatesByProjectID(projectID int64) ([]*EmailTemplate, error) {
	query := `
		SELECT
		  email_template_id,
		  project_id,
		  name,
		  is_active,
		  subject,
		  content,
		  created_at,
		  updated_at
		FROM
		  email_templates
		WHERE
		  project_id = $1
	`

	return scanEmailTemplates(query, projectID)
}

// getLastEmailTemplateByProjectID returns all email_templates in the database based
// on a given expression
func getActiveEmailTemplateByProjectID(projectID int64) (*EmailTemplate, error) {
	query := `
		SELECT
		  email_template_id,
		  project_id,
		  name,
		  is_active,
		  subject,
		  content,
		  created_at,
		  updated_at
		FROM
		  email_templates
		WHERE
		  project_id = $1
		  AND is_active = true
	`

	return scanEmailTemplate(query, projectID)
}

// updateEmailTemplate updates a single email_templates row in the database
func updateEmailTemplate(et *EmailTemplate) error {
	query := `
		UPDATE
		  email_templates
		SET
		  name=$1,
		  is_active=$2,
		  subject=$3,
		  content=$4
		WHERE
		  email_template_id = $5
	`

	if err := database.Exec(query, et.Name, et.IsActive, et.Subject, et.Content, et.ID); err != nil {
		return fmt.Errorf("failed updating email_template: %v", err)
	}

	return nil
}

// deleteEmailTemplate deletes a single email_templates row from database
func deleteEmailTemplate(id int64) error {
	query := `DELETE FROM email_templates WHERE email_template_id=$1`

	if err := database.Exec(query, id); err != nil {
		return fmt.Errorf("Failed deleting email_templates row: %v", err)
	}

	return nil
}

// scanEmailTemplate returns a single email_template based on the given query
func scanEmailTemplate(query string, params ...interface{}) (*EmailTemplate, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	row := db.QueryRow(query, params...)
	et := New()

	if err := row.Scan(&et.ID, &et.ProjectID, &et.Name, &et.IsActive, &et.Subject, &et.Content, &et.CreatedAt, &et.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed scanning email_templates row: %v", err)
		}

		return nil, nil
	}

	return et, nil
}

// scanEmailTemplates returns multiple email_templates based on the given query
func scanEmailTemplates(query string, params ...interface{}) ([]*EmailTemplate, error) {
	var ets []*EmailTemplate

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(query, params...)
	if err != nil {
		return ets, fmt.Errorf("Failed executing `%s`: %v", query, err)
	}

	defer rows.Close()

	for rows.Next() {
		et := New()

		if err := rows.Scan(&et.ID, &et.ProjectID, &et.Name, &et.IsActive, &et.Subject, &et.Content, &et.CreatedAt, &et.UpdatedAt); err != nil {
			return ets, fmt.Errorf("Failed scanning email_templates row: %v", err)
		}

		ets = append(ets, et)
	}

	return ets, nil
}
