package project

import (
	"fmt"
	"strconv"
)

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
