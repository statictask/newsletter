package task

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/database"
)

// insertTask inserts a Task in the database
func insertTask(t *Task) error {
	query := `
		INSERT INTO tasks (
		  pipeline_id,
		  task_type
	  	)
		VALUES (
		  $1,
		  $2
		)
		RETURNING
		  task_id,
		  pipeline_id,
		  task_type,
		  task_status,
		  created_at,
		  updated_at
	`

	savedTask, err := scanTask(query, t.PipelineID, t.Type)
	if err != nil {
		return err
	}

	*t = *savedTask

	return nil
}

// getTask return a single row that matches a given expression
func getTaskByTypeAndPipelineID(taskType string, pipelineID int64) (*Task, error) {
	query := `
		SELECT
		  task_id,
		  pipeline_id,
		  task_type,
		  task_status,
		  created_at,
		  updated_at
		FROM
		  tasks
		WHERE
		  task_type=$1 AND
		  pipeline_id=$2
	`

	return scanTask(query, taskType, pipelineID)
}

// getTasksByPipelineID returns all tasks in the database based on the pipeline ID
func getTasksByPipelineID(pipelineID int64) ([]*Task, error) {
	query := `
		SELECT
		  task_id,
		  pipeline_id,
		  task_type,
		  task_status,
		  created_at,
		  updated_at
		FROM
		  tasks
		WHERE
		  pipeline_id=$1
	`

	return scanTasks(query, pipelineID)
}

// getTasksByTypeAndStatus returns all tasks in the database based on the pipeline ID
func getTasksByTypeAndStatus(taskType, taskStatus string) ([]*Task, error) {
	query := `
		SELECT
		  task_id,
		  pipeline_id,
		  task_type,
		  task_status,
		  created_at,
		  updated_at
		FROM
		  tasks
		WHERE
		  task_type=$1 AND
		  task_status=$2
	`

	return scanTasks(query, taskType, taskStatus)
}

// updateTask updates a Task in the database
func updateTask(t *Task) error {
	// only allows update to the task_status field, the other fields are immutable
	query := `UPDATE tasks SET task_status=$1 WHERE task_id=$2`

	if err := database.Exec(query, t.Status, t.ID); err != nil {
		return fmt.Errorf("failed updating task: %v", err)
	}

	return nil
}

// scanTask returns a single task that matches the given query
func scanTask(query string, params ...interface{}) (*Task, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	row := db.QueryRow(query, params...)
	t := &Task{}

	if err := row.Scan(&t.ID, &t.PipelineID, &t.Type, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("unable to scan task row: %v", err)
		}

		return nil, nil
	}

	return t, nil
}

// scanTasks returns multiple tasks that match the given query
func scanTasks(query string, params ...interface{}) ([]*Task, error) {
	var ts []*Task

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(query, params...)
	if err != nil {
		return ts, fmt.Errorf("unable to execute `%s`: %v", query, err)
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
