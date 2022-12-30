package feedreader

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
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

func New(url string) *FeedReader {
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
