package postitem

import (
	"fmt"
	"time"
)

type PostItem struct {
	ID         int64
	PostID 	   int64
	Title      string
	Link       string
	Content    string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func New() *PostItem {
	return &PostItem{}
}

// Create the PostItem in the database
func (p *PostItem) Create() error {
	if err := insertPostItem(p); err != nil {
		return fmt.Errorf("unable to create post_item: %v", err)
	}

	return nil
}
