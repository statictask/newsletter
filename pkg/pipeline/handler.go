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

	sqlStatement := "SELECT pipeline_id,project_id,created_at,updated_at FROM pipelines ORDER BY created_at DESC LIMIT 1"

	row := db.QueryRow(sqlStatement)
	p := &Pipeline{}
	err = row.Scan(&p.ID, &p.ProjectID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("unable to scan pipeline row: %v", err)
		}

		return nil, nil

	}

	return p, nil
}

// getPipelines returns all pipelines in the database based
// on a given expression
func getPipelinesWhere(expression string) ([]*Pipeline, error) {
	var ps []*Pipeline

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT pipeline_id,project_id,created_at,updated_at FROM pipelines"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

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

// deletePipeline deletes a Pipeline record from database
func deletePipeline(pipelineID int64) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `DELETE FROM pipelines WHERE pipeline_id=$1`

	res, err := db.Exec(sqlStatement, pipelineID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with pipeline_id `%d`: %v",
			sqlStatement, pipelineID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info(
		"pipeline rows deleted",
		zap.Int64("total", rowsAffected),
		zap.Int64("pipeline_id", pipelineID),
	)

	return nil
}

// insertTask inserts a Task in the database
func insertTask(t *Task) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `INSERT INTO tasks (pipeline_id,task_type) VALUES ($1,$2) RETURNING task_id,task_status,created_at,updated_at`

	if err := db.QueryRow(
		sqlStatement,
		t.PipelineID,
		t.Type,
	).Scan(&t.ID, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
		log.L.Fatal("unable to execute the query", zap.Error(err))
	}

	log.L.Info(
		"successfully created task record",
		zap.Int64("pipeline_id", t.PipelineID),
		zap.Int64("task_id", t.ID),
	)

	return nil
}

// getTask return a single row that matches a given expression
func getTaskWhere(expression string) (*Task, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT task_id,pipeline_id,task_type,task_status,created_at,updated_at FROM tasks"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	row := db.QueryRow(sqlStatement)
	t := &Task{}

	if err := row.Scan(&t.ID, &t.PipelineID, &t.Type, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("unable to scan task row: %v", err)
		}

		return nil, nil
	}

	return t, nil
}

// getTasks returns all tasks in the database based
// on a given expression
func getTasksWhere(expression string) ([]*Task, error) {
	var ts []*Task

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT task_id,pipeline_id,task_type,task_status,created_at,updated_at FROM tasks"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return ts, fmt.Errorf("unable to execute `%s`: %v", sqlStatement, err)
	}

	defer rows.Close()

	for rows.Next() {
		t := NewTask()

		if err := rows.Scan(&t.ID, &t.PipelineID, &t.Type, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return ts, fmt.Errorf("unable to scan task row: %v", err)
		}

		ts = append(ts, t)
	}

	return ts, nil

}

// updateTask updates a Task in the database
func updateTask(t *Task) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	// only allows update to the task_status field, the other fields are immutable
	sqlStatement := `UPDATE tasks SET task_status WHERE task_id=$3`

	res, err := db.Exec(sqlStatement, t.Status, t.ID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with task_id `%d`: %v",
			sqlStatement, t.ID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info("tasks rows updated", zap.Int64("total", rowsAffected), zap.Int64("task_id", t.ID))

	return nil
}

// deleteTask deletes a Task record from database
func deleteTask(taskID int64) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `DELETE FROM tasks WHERE task_id=$1`

	res, err := db.Exec(sqlStatement, taskID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with task_id `%d`: %v",
			sqlStatement, taskID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info(
		"task rows deleted",
		zap.Int64("total", rowsAffected),
		zap.Int64("task_id", taskID),
	)

	return nil
}
