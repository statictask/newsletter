package task

// PipelineTasks is the entity used for lazy controlling
// interactions with many PipelineTasks in the database
type PipelineTasks struct {
	pipelineID int64
}

// NewPipelineTasks returns a PipelineTasks controller
func NewPipelineTasks(pipelineID int64) *PipelineTasks {
	return &PipelineTasks{pipelineID}
}

// All returns all the ps registered in the database
// for a given project
func (ts *PipelineTasks) All() ([]*Task, error) {
	return getTasksByPipelineID(ts.pipelineID)
}

// Add creates a new entry in the pipeline's tasks
func (ts *PipelineTasks) Create(taskType TaskType) (*Task, error) {
	// make sure the Task has the correct ProjectID
	t := NewTask()
	t.PipelineID = ts.pipelineID
	t.Type = taskType

	if err := t.Create(); err != nil {
		return nil, err
	}

	return t, nil
}

// GetByType return a single task of a given type
func (ts *PipelineTasks) GetByType(taskType TaskType) (*Task, error) {
	return getTaskByTypeAndPipelineID(string(taskType), ts.pipelineID)
}
