package watcher

import (
	"context"

	"github.com/mmcdole/gofeed"
)

type FeedReader struct {
	Url string
}

func NewFeedReader(url string) *FeedReader {
	return &FeedReader{url}
}

func (f *FeedReader) Run(ctx context.Context) (string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(f.Url, ctx)
	if err != nil {
		return "", err
	}

	return feed.Title, nil
}
