package feedreader

import "time"

type Event struct {
	ID           int64
	FeedReaderID int64
	ContentHash  string
	CreatedAt    time.Time
}

func NewEvent() *Event {
	return &Event{}
}

// Create records a new Event in the database
func (e *Event) Create() error {
	if err := insertEvent(e); err != nil {
		return fmt.Errorf("unable to create feed reader event: %v", err)
	}

	return nil
}

// Delete an event from the database
func (e *Event) Delete() error {
	if err := deleteEvent(e.ID, e.FeedReaderID); err != nil {
		return fmt.Errorf("unable to delete feed reader event: %v", err)
	}

	return nil
}
