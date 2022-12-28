package pipeline

import (
	"fmt"
	"time"

	"github.com/statictask/newsletter/pkg/post"
)

type Pipeline struct {
	ID        int64
	ProjectID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New() *Pipeline {
	return &Pipeline{}
}

// Create the Pipeline in the database
func (p *Pipeline) Create() error {
	if err := insertPipeline(p); err != nil {
		return fmt.Errorf("unable to create pipeline: %v", err)
	}

	return nil
}

// Delete the Pipeline from the database
func (p *Pipeline) Delete() error {
	if err := deletePipeline(p.ID); err != nil {
		return fmt.Errorf("unable to delete pipeline: %v", err)
	}

	return nil
}

// Tasks returns a lazy interface for interacting with Pipelines's Task objects
func (p *Pipeline) Tasks() *Tasks {
	return NewTasks(p.ID)
}

// Posts returns a lazy interface for interacting with Pipelines's Post objects
func (p *Pipeline) Posts() *post.Posts {
	return post.NewPosts(p.ID)
}

// IsFinished checks if the pipeline is still running by querying the state of
// inner tasks
func (p *Pipeline) IsFinished() (bool, error) {
	is_scraped := false
	is_published := false

	tasks, err := p.Tasks().All()
	if err != nil {
		return false, err
	}

	for _, t := range tasks {
		if t.IsScrape() && t.IsFinished() {
			is_scraped = true
			continue
		}

		if t.IsPublish() && t.IsFinished() {
			is_published = true
		}
	}

	return is_scraped && is_published, nil
}
