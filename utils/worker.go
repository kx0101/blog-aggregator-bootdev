package utils

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
}

func fetchRSSFeed(url string) (*RSS, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %v", err)
	}

	return &rss, nil
}

func FeedWorker(dbQueries *database.Queries, interval time.Duration, batchSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		feeds, err := dbQueries.GetNextFeedsToFetch(ctx, 10)
		if err != nil {
			fmt.Printf("Error fetching RSS feed %v: %s\n", feeds, err)
			return
		}

		var wg sync.WaitGroup
		wg.Add(len(feeds))

		for _, feed := range feeds {
			go func(feed database.Feed) {
				defer wg.Done()

				rss, err := fetchRSSFeed(feed.Url)
				if err != nil {
					fmt.Printf("Error fetching RSS feed %s: %s\n", feed.Url, err)
					return
				}

				for _, item := range rss.Channel.Items {
					log.Printf("Feed: %s - Post: %s", rss.Channel.Title, item.Title)
				}

				if err := dbQueries.MarkFeedFetched(ctx, feed.ID); err != nil {
					log.Printf("Error marking feed %d as fetched: %s\n", feed.ID, err)
				}
			}(feed)
		}

		wg.Wait()

		time.Sleep(interval)
	}
}
