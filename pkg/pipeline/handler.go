package pipeline

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/database"
)

// insertPipeline inserts a Pipeline in the database
func insertPipeline(p *Pipeline) error {
	query := `
		INSERT INTO pipelines
		  (project_id)
		VALUES
		  ($1)
		RETURNING
		  pipeline_id,
		  project_id,
		  created_at,
		  updated_at
	`

	_, err := scanPipeline(query, p.ProjectID)
	return err
}

// getLastPipeline return a single row that matches the last pipeline for a given project_id
func getLastPipeline(projectID int64) (*Pipeline, error) {
	query := `
		SELECT
		  pipeline_id,
		  project_id,
		  created_at,
		  updated_at
		FROM
		  pipelines
		ORDER BY
		  created_at
		DESC
		LIMIT 1
	`
	return scanPipeline(query)
}

// getPipelinesByProjectID returns all pipelines in the database based
// on a given expression
func getPipelinesByProjectID(projectID int64) ([]*Pipeline, error) {
	template := `
		SELECT 
		  pipeline_id,project_id,created_at,updated_at
		FROM
		  pipelines
		WHERE
		  project_id = %d
	`

	query := fmt.Sprintf(template, projectID)
	return scanPipelines(query)
}

// scanPipelines returns multiple pipelines that match the given query
func scanPipelines(query string) ([]*Pipeline, error) {
	var ps []*Pipeline

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return ps, fmt.Errorf("unable to execute `%s`: %v", query, err)
	}

	defer rows.Close()

	for rows.Next() {
		p := New()

		if err := rows.Scan(&p.ID, &p.ProjectID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return ps, fmt.Errorf("unable to scan pipeline row: %v", err)
		}

		ps = append(ps, p)
	}

	return ps, nil
}

// scanPipeline returns a single pipeline that matches the given query
func scanPipeline(query string, params ...interface{}) (*Pipeline, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	row := db.QueryRow(query, params...)

	p := &Pipeline{}
	if err := row.Scan(&p.ID, &p.ProjectID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("unable to scan pipeline row: %v", err)
		}

		return nil, nil
	}

	return p, nil
}
