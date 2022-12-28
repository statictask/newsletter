package post

import (
	"fmt"
	"time"
)

type Post struct {
	ID         int64
	PipelineID int64
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func New() *Post {
	return &Post{}
}

// Create the Post in the database
func (p *Post) Create() error {
	if err := insertPost(p); err != nil {
		return fmt.Errorf("unable to create post: %v", err)
	}

	return nil
}

// Update the Post record in the database
func (p *Post) Update() error {
	if err := updatePost(p); err != nil {
		return fmt.Errorf("unable to update post: %v", err)
	}

	return nil
}

// Delete the Task from the database
func (p *Post) Delete() error {
	if err := deletePost(p.ID); err != nil {
		return fmt.Errorf("unable to delete post: %v", err)
	}

	return nil
}
