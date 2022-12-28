package project

import (
	"fmt"
	"time"

	"github.com/statictask/newsletter/pkg/pipeline"
	"github.com/statictask/newsletter/pkg/subscription"
)

type Project struct {
	ID        int64     `json:"project_id"`
	Name      string    `json:"name"`
	FeedURL   string    `json:"feed_url"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

// Subscriptions return a lazy inteface for interacting with project's subscriptions
func (p *Project) Subscriptions() *subscription.Subscriptions {
	return subscription.NewSubscriptions(p.ID)
}

// Pipelines return a lazy inteface for interacting with project's pipelines
func (p *Project) Pipelines() *pipeline.Pipelines {
	return pipeline.NewPipelines(p.ID)
}
