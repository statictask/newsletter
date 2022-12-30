package task

import (
	"fmt"
	"time"
)

type TaskType string
type TaskStatus string

const (
	Scrape  TaskType = "Scrape"
	Publish TaskType = "Publish"

	Waiting  TaskStatus = "Waiting"
	Ready    TaskStatus = "Ready"
	Running  TaskStatus = "Running"
	Finished TaskStatus = "Finished"
	Failed   TaskStatus = "Failed"
	Aborted  TaskStatus = "Aborted"
)

var (
	TaskTypes    []TaskType   = []TaskType{Scrape, Publish}
	TaskStatuses []TaskStatus = []TaskStatus{Waiting, Ready, Running, Finished, Failed, Aborted}
)

type Task struct {
	ID         int64
	PipelineID int64
	Type       TaskType
	Status     TaskStatus
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func NewTask() *Task {
	return &Task{}
}

// Create the Task in the database
func (t *Task) Create() error {
	if err := insertTask(t); err != nil {
		return fmt.Errorf("unable to create task: %v", err)
	}

	return nil
}

// Update the Task in the database
func (t *Task) Update() error {
	if err := updateTask(t); err != nil {
		return fmt.Errorf("unable to update task: %v", err)
	}

	return nil
}

// IsFinished says if the status of the task is Finished or not
func (t *Task) IsFinished() bool {
	return t.Status == Finished
}

// IsScrape says if the type of the task is Scrape or not
func (t *Task) IsScrape() bool {
	return t.Type == Scrape
}

// IsScrape says if the type of the task is Publish or not
func (t *Task) IsPublish() bool {
	return t.Type == Publish
}
