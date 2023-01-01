package project

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/database"
)

// insertProject inserts a project in the database
func insertProject(p *Project) error {
	query := `
		INSERT INTO projects (
		  name,
		  feed_url
	  	)
		VALUES (
		  $1,
		  $2
		)
		RETURNING
		  project_id,
		  name,
		  feed_url,
		  is_enabled,
		  created_at,
		  updated_at
	`

	savedProject, err := scanProject(query, p.Name, p.FeedURL)
	if err != nil {
		return err
	}

	*p = *savedProject

	return nil
}

// updateProject updates a project in the database
func updateProject(p *Project) error {
	query := `
		UPDATE
		  projects
		SET
		  name=$1,
		  feed_url=$2,
		  is_enabled=$3
		WHERE
		  project_id=$4
	`

	if err := database.Exec(query, p.Name, p.FeedURL, p.IsEnabled, p.ID); err != nil {
		return fmt.Errorf("failed updating project: %v", err)
	}

	return nil
}

// deleteProject deletes a project from database
func deleteProject(projectID int64) error {
	query := `DELETE FROM projects WHERE project_id=$1`

	if err := database.Exec(query, projectID); err != nil {
		return fmt.Errorf("failed deleting project: %v", err)
	}

	return nil
}

// getEnabledProjects returns projects that match is_enabled = true
func getEnabledProjects() ([]*Project, error) {
	query := `
		SELECT
		  project_id,
		  name,
		  feed_url,
		  is_enabled,
		  created_at,
		  updated_at
		FROM
		  projects
		WHERE
		  is_enabled = true 
	`

	return scanProjects(query)
}

// getProjectByTask returns the project associated with the given task id
func getProjectByTaskID(taskID int64) (*Project, error) {
	query := `
		SELECT
		  pr.project_id,
		  pr.name,
		  pr.feed_url,
		  pr.is_enabled,
		  pr.created_at,
		  pr.updated_at
		FROM
		  projects AS pr
		JOIN pipelines AS pi
		  ON pr.project_id = pi.project_id
		FULL OUTER JOIN tasks AS ta
		  ON ta.pipeline_id = pi.pipeline_id
		WHERE
		  ta.task_id=$1
	`
	return scanProject(query, taskID)
}

// getProjectByID returns a single project based on its ID
func getProjectByID(projectID int64) (*Project, error) {
	query := `
		SELECT
		  project_id,
		  name,
		  feed_url,
		  is_enabled,
		  created_at,
		  updated_at
		FROM
		  projects
		WHERE
		  project_id = $1
	`

	return scanProject(query, projectID)
}

// scanProject returns a single project based on the given query
func scanProject(query string, params ...interface{}) (*Project, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	row := db.QueryRow(query, params...)
	p := New()

	if err := row.Scan(&p.ID, &p.Name, &p.FeedURL, &p.IsEnabled, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("unable to scan project row: %v", err)
		}

		return nil, nil
	}

	return p, nil
}

// scanProjects returns multiple projects that match the given query
func scanProjects(query string, params ...interface{}) ([]*Project, error) {
	var projects []*Project

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return projects, fmt.Errorf("unable to execute `%s`: %v", query, err)
	}

	defer rows.Close()

	for rows.Next() {
		p := New()

		if err := rows.Scan(&p.ID, &p.Name, &p.FeedURL, &p.IsEnabled, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return projects, fmt.Errorf("unable to scan a project row: %v", err)
		}

		projects = append(projects, p)
	}

	return projects, nil
}
