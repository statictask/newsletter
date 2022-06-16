package project

import (
	"fmt"

	"github.com/statictask/newsletter/pkg/subscription"
)

type Project struct {
	ID     int64  `json:"project_id"`
	Domain string `json:"domain"`
}

// New returns an empty Project
func New() *Project {
	return &Project{}
}

// Create the project in the database
func (p *Project) Create() error {
	if err := insertProject(p); err != nil {
		return fmt.Errorf("unable to create project: %v", err)
	}

	return nil
}

// Update the project in the database
func (p *Project) Update() error {
	if err := updateProject(p); err != nil {
		return fmt.Errorf("unable to update project: %v", err)
	}

	return nil
}

// Delete the project from the database
func (p *Project) Delete() error {
	if err := deleteProject(p.ID); err != nil {
		return fmt.Errorf("unable to delete project: %v", err)
	}

	return nil
}

// Subscriptions return an inteface for interact with Project's subscriptions
func (p *Project) Subscriptions() *subscription.Subscriptions {
	return subscription.NewSubscriptions(p.ID)
}
