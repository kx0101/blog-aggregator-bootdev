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

	"github.com/google/uuid"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Items       []Item    `xml:"item"`
	PublishedAt time.Time `xml:"pubDate"`
}

type Item struct {
	Title string `xml:"title"`
}

func FeedWorker(dbQueries *database.Queries, interval time.Duration, batchSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		feeds, err := dbQueries.GetNextFeedsToFetch(ctx, int32(batchSize))
		if err != nil {
			fmt.Printf("Error fetching RSS feed %v: %s\n", feeds, err)
			cancel()
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

				id := uuid.New()
				now := time.Now()
				_, err = dbQueries.CreatePost(ctx, database.CreatePostParams{
					ID:          id,
					CreatedAt:   now,
					UpdatedAt:   now,
					Title:       rss.Channel.Title,
					Url:         feed.Url,
					PublishedAt: feed.CreatedAt,
					FeedID:      feed.ID,
				})
				if err != nil {
					fmt.Printf("Error creating post: %s", err)
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
		cancel()
	}
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
