package project

import (
	"fmt"
	"strconv"

	"github.com/statictask/newsletter/internal/database"
	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

// insertProject inserts a project in the database
func insertProject(p *Project) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `INSERT INTO projects (name,feed_url) VALUES ($1,$2) RETURNING project_id,is_enabled,created_at,updated_at`

	if err := db.QueryRow(
		sqlStatement,
		p.Name,
		p.FeedURL,
	).Scan(&p.ID, &p.IsEnabled, &p.CreatedAt, &p.UpdatedAt); err != nil {
		log.L.Fatal("unable to execute the query", zap.Error(err))
	}

	log.L.Info(
		"created project record",
		zap.Int64("project_id", p.ID),
		zap.String("name", p.Name),
	)

	return nil
}

// getProjectWhere return a single row that matches a given expression
func getProjectWhere(expression string) (*Project, error) {
	query := "SELECT project_id,name,feed_url,is_enabled,created_at,updated_at FROM projects"

	if expression != "" {
		query = fmt.Sprintf("%s WHERE %s", query, expression)
	}

	return scanProject(query)
}

// getProjects returns all projects in the database
func getProjectsWhere(expression string) ([]*Project, error) {
	var projects []*Project

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT project_id,name,feed_url,is_enabled,created_at,updated_at FROM projects"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return projects, fmt.Errorf("unable to execute `%s`: %v", sqlStatement, err)
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

// updateProject updates a project in the database
func updateProject(p *Project) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `UPDATE projects SET name=$1,feed_url=$2,is_enabled=$3 WHERE project_id=$4`

	res, err := db.Exec(sqlStatement, p.Name, p.FeedURL, p.IsEnabled, p.ID)
	if err != nil {
		return fmt.Errorf("unable to execute `%s`: %v", sqlStatement, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info(
		"project rows updated",
		zap.Int64("total", rowsAffected),
		zap.Int64("project_id", p.ID),
	)

	return nil
}

// deleteProject deletes a project from database
func deleteProject(projectID int64) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `DELETE FROM projects WHERE project_id=$1`

	res, err := db.Exec(sqlStatement, projectID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with project_id `%v`: %v",
			sqlStatement, projectID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info("project rows deleted", zap.Int64("total", rowsAffected), zap.Int64("project_id", projectID))

	return nil
}

// getProjectByTask
func getProjectByTaskID(taskID int64) (*Project, error) {
	query := fmt.Sprintf(
		`SELECT pr.project_id, pr.name, pr.feed_url, pr.is_enabled, pr.created_at, pr.updated_at
			FROM projects AS pr
			JOIN pipelines AS pi
				ON pr.project_id = pi.project_id
			FULL OUTER JOIN tasks AS ta
				ON ta.pipeline_id = pi.pipeline_id
			WHERE ta.task_id=%d`,
		taskID,
	)

	return scanProject(query)
}

// scanProject returns a single project based on the given query
func scanProject(query string) (*Project, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	row := db.QueryRow(query)
	p := New()

	if err := row.Scan(&p.ID, &p.Name, &p.FeedURL, &p.IsEnabled, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, fmt.Errorf("unable to scan a project row: %v", err)
	}

	return p, nil
}

// loadProject is a helper function that receives an string with the
// project_id and returns a project instance loaded from db
func loadProject(projectID string) (*Project, error) {
	id, err := castID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed casting id: %v", err)
	}

	projects := NewProjects()

	project, err := projects.Get(int64(id))
	if err != nil {
		return nil, fmt.Errorf("failed retrieving project: %v", err)
	}

	return project, nil
}

// castID converts a string ID to an int64 ID
func castID(strID string) (int64, error) {
	id, err := strconv.Atoi(strID)
	if err != nil {
		return -1, fmt.Errorf("unable to parse project_id into int: %v", err)
	}

	return int64(id), nil
}
