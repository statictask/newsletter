package feedreader

import "fmt"

// FeedReaders is the entity used for lazy controlling
// interactions with many FeedReaders in the database
type FeedReaders struct {
	projectID int64
}

// NewFeedReaders returns a FeedReaders controller
func NewFeedReaders(projectID int64) *FeedReaders {
	return &FeedReaders{projectID}
}

// All returns all the frs registered in the database
// for a given project
func (frs *FeedReaders) All() ([]*FeedReader, error) {
	exp := fmt.Sprintf("project_id=%d", frs.projectID)

	frArray, err := getFeedReadersWhere(exp)
	if err != nil {
		return []*FeedReader{}, fmt.Errorf("unable to get feed readers: %v", err)
	}

	return frArray, nil
}

// Get returns a single FeedReader according to its ID
func (frs *FeedReaders) Get(feedReaderID int64) (*FeedReader, error) {
	exp := fmt.Sprintf("feed_reader_id=%d AND project_id=%d", feedReaderID, frs.projectID)

	return getFeedReaderWhere(exp)
}

// Where return many FeedReaders according to an expression
func (frs *FeedReaders) Where(exp string) ([]*FeedReader, error) {
	exp = fmt.Sprintf("%s AND project_id=%d", exp, frs.projectID)

	frArray, err := getFeedReadersWhere(exp)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed readers: %v", err)
	}

	return frArray, nil
}

// Delete deletes a FeedReader based on its ID
func (frs *FeedReaders) Delete(feedReaderID int64) error {
	if err := deleteFeedReader(feedReaderID, frs.projectID); err != nil {
		return fmt.Errorf("unable to delete feed reader: %v", err)
	}

	return nil
}

// Add creates a new entry in the project's frs
// the function creates a new feedReader entry in the database
func (frs *FeedReaders) Add(fr *FeedReader) error {
	// make sure the FeedReader has the correct ProjectID
	fr.ProjectID = frs.projectID

	if err := fr.Create(); err != nil {
		return err
	}

	return nil
}
