package watcher

import (
	"context"
	"time"

	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

type FeedWatcher struct {
	reader *FeedReader
}

func NewFeedWatcher(url string) *FeedWatcher {
	reader := NewFeedReader(url)
	return &FeedWatcher{reader}
}

func (fw *FeedWatcher) Run(message chan string) chan WatcherSignal {
	stop := make(chan WatcherSignal)

	fw.watch(stop, message)
	return stop
}

func (fw *FeedWatcher) watch(stop chan WatcherSignal, message chan string) {
	go func() {
		for {
			select {
			case <-stop:
				log.L.Info("stopping feed watcher", zap.String("url", fw.reader.Url))
				return
			default:
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				feedTitle, err := fw.reader.Run(ctx)
				if err != nil {
					log.L.Error("feed watcher failed reading feed", zap.String("url", fw.reader.Url))
				} else {
					message <- feedTitle
				}

				cancel()
				time.Sleep(20 * time.Second)
			}
		}
	}()
}
