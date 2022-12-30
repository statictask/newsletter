package project

import "fmt"

// Projects is the entity used for controlling
// interactions with many projects in the database
type Projects struct{}

// NewProjects returns a Projects controller
func NewProjects() *Projects {
	return &Projects{}
}

// All returns all the projects registered in the database
func (pp *Projects) All() ([]*Project, error) {
	projects, err := getProjectsWhere("")
	if err != nil {
		return projects, fmt.Errorf("unable to get projects: %v", err)
	}

	return projects, nil
}

// Get returns a single project according to its ID
func (pp *Projects) Get(projectID int64) (*Project, error) {
	exp := fmt.Sprintf("project_id=%v", projectID)

	return getProjectWhere(exp)
}

// Where return many projects according to a map of attrs
func (pp *Projects) Where(exp string) ([]*Project, error) {
	projects, err := getProjectsWhere(exp)
	if err != nil {
		return nil, fmt.Errorf("unable to get projects: %v", err)
	}

	return projects, nil
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
	return pp.Where("is_enabled = true")
}
