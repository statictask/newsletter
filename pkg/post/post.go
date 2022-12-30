package post

import (
	"fmt"
	"time"
)

type Post struct {
	ID         int64
	PipelineID int64
	Content    string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
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
