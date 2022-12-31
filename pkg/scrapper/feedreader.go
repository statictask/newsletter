package scrapper

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

type FeedReader struct {
	FeedURL string
}

type FeedItem struct {
	Title       string
	Description string
	Content     string
	Link        string
	PubDate     *time.Time
}

func NewFeedReader(url string) *FeedReader {
	return &FeedReader{url}
}

func (fr *FeedReader) ReadFrom(ctx context.Context, dateLimit *time.Time) ([]*FeedItem, error) {
	items := []*FeedItem{}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(fr.FeedURL, ctx)
	if err != nil {
		return items, err
	}

	for _, i := range feed.Items {
		log.L.Info("processing feed item", zap.String("item_guid", i.GUID))
		if i.PublishedParsed.After(*dateLimit) {
			item := &FeedItem{
				Title:       i.Title,
				Description: i.Description,
				Content:     i.Content,
				Link:        i.Link,
				PubDate:     i.PublishedParsed,
			}

			items = append(items, item)
		}
	}

	return items, nil
}

func (fi *FeedItem) GetTitle() string {
	return fi.Title
}

func (fi *FeedItem) GetContent() string {
	if fi.Content == "" {
		return fi.Description
	}

	return fi.Content
}

func (fi *FeedItem) GetLink() string {
	return fi.Link
}
