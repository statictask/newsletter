package post

import (
	"fmt"
	"time"

	"github.com/statictask/newsletter/pkg/postitem"
)

type Post struct {
	ID         int64
	PipelineID int64
	Title      string
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

// PostItems returns a lazy interface for interacting with
// post_items related to this post
func (p *Post) PostItems() *postitem.PostPostItems {
	return postitem.NewPostPostItems(p.ID)
}

// GetPlainTextContent returns a content of the post
// aggreggating all the post items in a single Plain 
// Text document
func (p *Post) GetPlainTextContent() (string, error) {
	items, err := p.PostItems().All()
	if err != nil {
		return "", err
	}

	builder := NewContentBuilder(p.Title, items)
	return builder.BuildPlainText()
}

// GetPlainTextContent returns a content of the post
// aggreggating all the post items in a single HTML
// document
func (p *Post) GetHTMLContent() (string, error) {
	items, err := p.PostItems().All()
	if err != nil {
		return "", err
	}

	builder := NewContentBuilder(p.Title, items)
	return builder.BuildHTML()
}

// GetSubject returns the Title of the post
func (p *Post) GetSubject() string {
	return p.Title
}
