package task

// Tasks is the entity used for lazy controlling
type Tasks struct{}

// NewTasks returns a Tasks controller
func NewTasks() *Tasks {
	return &Tasks{}
}

// Filter selects only tasks that matches both task_type and task_status
func (ts *Tasks) Filter(taskType TaskType, taskStatus TaskStatus) ([]*Task, error) {
	return getTasksByTypeAndStatus(string(taskType), string(taskStatus))
}
