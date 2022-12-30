package task

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/database"
	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

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
func getTaskByTypeAndPipelineID(taskType string, pipelineID int64) (*Task, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := fmt.Sprintf(
		"SELECT task_id,pipeline_id,task_type,task_status,created_at,updated_at FROM tasks WHERE task_type='%v' AND pipeline_id=%d",
		taskType,
		pipelineID,
	)

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

// getTasksByPipelineID returns all tasks in the database based on the pipeline ID
func getTasksByPipelineID(pipelineID int64) ([]*Task, error) {
	exp := fmt.Sprintf("pipeline_id='%d'", pipelineID)
	return getTasksWhere(exp)
}

// getTasksByTypeAndStatus returns all tasks in the database based on the pipeline ID
func getTasksByTypeAndStatus(taskType, taskStatus string) ([]*Task, error) {
	exp := fmt.Sprintf("task_type='%s' AND task_status='%s'", taskType, taskStatus)
	return getTasksWhere(exp)
}

// getTasksWhere returns all tasks in the database based on the pipeline ID
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
	sqlStatement := `UPDATE tasks SET task_status=$1 WHERE task_id=$2`

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
