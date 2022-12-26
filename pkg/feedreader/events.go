package feedreader

import "fmt"

type Events struct {
	feedReaderID int64
}

func NewEvents(feedReaderID int64) *Events {
	return &Events{feedReaderID}
}

func (es *Events) Last() (*Event, error) {
	exp := fmt.Sprintf("feed_reader_id=%d", es.feedReaderID)

	e, err := getLastEventWhere(exp)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get the last event for the feed reader `%d`: %v",
			es.feedReaderID,
			err,
		)
	}

	return e, nil
}

func (es *Events) Add(e *Event) error {
	// make sure the Event has the correct FeedReaderID
	e.FeedReaderID = es.feedReaderID

	if err := e.Create(); err != nil {
		return err
	}

	return nil
}
