package pipeline

import "fmt"

// Tasks is the entity used for lazy controlling
// interactions with many Tasks in the database
type Tasks struct {
	pipelineID int64
}

// NewTasks returns a Tasks controller
func NewTasks(pipelineID int64) *Tasks {
	return &Tasks{pipelineID}
}

// All returns all the ps registered in the database
// for a given project
func (ts *Tasks) All() ([]*Task, error) {
	exp := fmt.Sprintf("pipeline_id=%d", ts.pipelineID)

	tArray, err := getTasksWhere(exp)
	if err != nil {
		return []*Task{}, fmt.Errorf("unable to get tasks: %v", err)
	}

	return tArray, nil
}

// Get returns a single Task according to its ID
func (ts *Tasks) Get(taskID int64) (*Task, error) {
	exp := fmt.Sprintf("task_id=%d AND pipeline_id=%d", taskID, ts.pipelineID)

	return getTaskWhere(exp)
}

// Where return many Tasks according to an expression
func (ts *Tasks) Where(exp string) ([]*Task, error) {
	exp = fmt.Sprintf("%s AND pipeline_id=%d", exp, ts.pipelineID)

	tArray, err := getTasksWhere(exp)
	if err != nil {
		return nil, fmt.Errorf("unable to get tasks: %v", err)
	}

	return tArray, nil
}

// Delete deletes a Task based on its ID
func (ts *Tasks) Delete(taskID int64) error {
	if err := deleteTask(taskID); err != nil {
		return fmt.Errorf("unable to delete task: %v", err)
	}

	return nil
}

// Add creates a new entry in the pipeline's tasks
func (ts *Tasks) Add(t *Task) error {
	// make sure the Task has the correct ProjectID
	t.PipelineID = ts.pipelineID

	if err := t.Create(); err != nil {
		return err
	}

	return nil
}
