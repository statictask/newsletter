package pipeline

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/database"
	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

// insertPipeline inserts a Pipeline in the database
func insertPipeline(p *Pipeline) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `INSERT INTO pipelines (project_id) VALUES ($1) RETURNING pipeline_id,created_at,updated_at`

	if err := db.QueryRow(
		sqlStatement,
		p.ProjectID,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		log.L.Fatal("unable to execute the query", zap.Error(err))
	}

	log.L.Info("successfully created pipeline record", zap.Int64("project_id", p.ProjectID), zap.Int64("pipeline_id", p.ID))

	return nil
}

// getPipeline return a single row that matches a given expression
func getPipelineWhere(expression string) (*Pipeline, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT pipeline_id,project_id,created_at,updated_at FROM pipelines"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	row := db.QueryRow(sqlStatement)
	p := &Pipeline{}

	if err := row.Scan(&p.ID, &p.ProjectID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, fmt.Errorf("unable to scan pipeline row: %v", err)
	}

	return p, nil
}


// getLastPipeline return a single row that matches the last pipeline for a given project_id
func getLastPipeline(projectID int64) (*Pipeline, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	query := "SELECT pipeline_id,project_id,created_at,updated_at FROM pipelines ORDER BY created_at DESC LIMIT 1"

	row := db.QueryRow(query)
	p := &Pipeline{}
	if err := row.Scan(&p.ID, &p.ProjectID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("unable to scan pipeline row: %v", err)
		}

		return nil, nil
	}

	return p, nil
}

// getPipelinesByProjectID returns all pipelines in the database based
// on a given expression
func getPipelinesByProjectID(projectID int64) ([]*Pipeline, error) {
	var ps []*Pipeline

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := fmt.Sprintf(
		"SELECT pipeline_id,project_id,created_at,updated_at FROM pipelines WHERE project_id = %d",
		projectID,
	)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return ps, fmt.Errorf("unable to execute `%s`: %v", sqlStatement, err)
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
