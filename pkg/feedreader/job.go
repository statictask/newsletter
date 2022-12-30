package feedreader

type FeedReaderJob struct {
	projectID int64
}

func NewFeedReaderJob(projectID int64) *FeedReaderJob {
	return &FeedReaderJob{projectID}
}

// Run checks if there are new items published in the feed since the last
// finished pipeline. If so, it'll get these items and create a post with
// the new content in the database
func (j *FeedReaderJob) Run() error {
	return nil
}
