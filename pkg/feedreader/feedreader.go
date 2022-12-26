package feedreader

import (
	"context"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedReader struct {
	ID        int64
	ProjectID int64
	FeedURL   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New() *FeedReader {
	return &FeedReader{}
}

// Create the FeedReader in the database
func (fr *FeedReader) Create() error {
	if err := insertFeedReader(fr); err != nil {
		return fmt.Errorf("unable to create feed reader: %v", err)
	}

	initialEvent := NewEvent()
	initialEvent.ContentHash = ""

	if err := fr.Events().Add(initialEvent) {
		return fmt.Errorf("failed creating initial feed reader event: %v", err)
	}

	return nil
}

// Update the FeedReader in the database
func (fr *FeedReader) Update() error {
	if err := updateFeedReader(fr); err != nil {
		return fmt.Errorf("unable to update feed reader: %v", err)
	}

	return nil
}

// Delete the FeedReader from the database
func (fr *FeedReader) Delete() error {
	if err := deleteFeedReader(fr.ID, fr.ProjectID); err != nil {
		return fmt.Errorf("unable to delete feed reader: %v", err)
	}

	return nil
}

func (fr *FeedReader) Read(ctx context.Context) (string, error) {
	// TODO: implement read based on events
	// TODO: check if it's possible to bring the watcher logic to this place
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(f.Url, ctx)
	if err != nil {
		return "", err
	}

	return feed.Title, nil
}

// Events returns a lazy interface for interacting with FeedReader's events
func (fr *FeedReader) Events() *Events {
	return NewEvents(fr.ID)
}
