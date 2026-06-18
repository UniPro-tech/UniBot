package util

import (
	"io"

	"github.com/mmcdole/gofeed"
)

const (
	MaxFeedSize = 10 * 1024 * 1024
	MaxRedirect = 5
)

func FetchFeed(feedURL string) (*gofeed.Feed, error) {
	body, err := HttpGet(feedURL)
	if err != nil || body == nil {
		return nil, err
	}
	parser := gofeed.NewParser()

	feed, err := parser.Parse(
		io.LimitReader(*body, MaxFeedSize),
	)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
