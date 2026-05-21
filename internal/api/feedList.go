package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	c := NewClient(1 * time.Minute)
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request %w", err)
	}
	req.Header.Set("User-Agent", "gator")
	req.Header.Set("Accept", "application/xml")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot receive response %w", err)
	}
	defer res.Body.Close()
	decoder := xml.NewDecoder(res.Body)
	var feed RSSFeed
	err = decoder.Decode(&feed)
	if err != nil {
		return nil, fmt.Errorf("issue when decoding response %w", err)
	}
	unescape(&feed)
	return &feed, nil
}

func unescape(rSSFeed *RSSFeed) {
	rSSFeed.Channel.Title = html.UnescapeString(rSSFeed.Channel.Title)
	rSSFeed.Channel.Description = html.UnescapeString(rSSFeed.Channel.Description)
	for i := range rSSFeed.Channel.Item {
		rSSFeed.Channel.Item[i].Title = html.UnescapeString(rSSFeed.Channel.Item[i].Title)
		rSSFeed.Channel.Item[i].Description = html.UnescapeString(rSSFeed.Channel.Item[i].Description)
	}
}
