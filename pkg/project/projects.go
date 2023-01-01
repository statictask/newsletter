package project

import "fmt"

// Projects is the entity used for controlling
// interactions with many projects in the database
type Projects struct{}

// NewProjects returns a Projects controller
func NewProjects() *Projects {
	return &Projects{}
}

// Get returns a single project according to its ID
func (pp *Projects) Get(projectID int64) (*Project, error) {
	return getProjectByID(projectID)
}

// Delete deletes a project based on its ID
func (pp *Projects) Delete(projectID int64) error {
	if err := deleteProject(projectID); err != nil {
		return fmt.Errorf("unable to delete project: %v", err)
	}

	return nil
}

// AllEnabled returns all the projects registered in the database that are enabled
func (pp *Projects) AllEnabled() ([]*Project, error) {
	return getEnabledProjects()
}

// GetByTaskID returns the project based on any related task
func (pp *Projects) GetByTaskID(taskID int64) (*Project, error) {
	return getProjectByTaskID(taskID)
}
