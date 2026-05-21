package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/debugger1709/blog-aggregator/internal/api"
	"github.com/debugger1709/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregate(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("a time duration is required")
	}
	delta, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("cannot parse %s to time duration", cmd.args[0])
	}
	ticker := time.NewTicker(delta)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s, user)
		if err != nil {
			return err
		}
	}
}

func handlerAllFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("issue when getting all feeds %w", err)
	}
	for _, feed := range feeds {
		fmt.Println(feed)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("a feed name and a url is required")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}
	f, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("issue when creating feed for user %s with feed name %s %w", feedName, feedURL, err)
	}
	fmt.Printf("feed %s was added by user %s\n", feedName, feedURL)
	fmt.Println(f)
	return createFeedFollow(user, f, s)
}

func scrapeFeeds(s *state, user database.User) error {
	f, err := s.db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("cannot fetch the next feed for user %s %w", user.Name, err)
	}
	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: time.Now().UTC(),
		ID:        f.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return fmt.Errorf("cannot mark feed %w", err)
	}
	rss, err := api.FetchFeed(context.Background(), f.Url)
	if err != nil {
		return fmt.Errorf("cannot fet rss feed from %s %w", f.Url, err)
	}
	return createPosts(s, rss, f)
}

func createPosts(s *state, rss *api.RSSFeed, f database.Feed) error {
	for i := range rss.Channel.Item {
		cur := rss.Channel.Item[i]
		params := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     cur.Title,
			Url:       cur.Link,
			Description: sql.NullString{
				String: cur.Description,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time:  convertTime(cur.PubDate),
				Valid: true,
			},
			FeedID: f.ID,
		}
		_, err := s.db.CreatePost(context.Background(), params)
		if err != nil {
			return fmt.Errorf("cannot create post %v %w", params, err)
		}
	}
	return nil
}

func convertTime(str string) time.Time {
	if t, err := time.Parse(time.RFC1123, str); err == nil {
		return t
	}
	if t, err := time.Parse(time.RFC1123Z, str); err == nil {
		return t
	}
	if t, err := time.Parse(time.RFC822, str); err == nil {
		return t
	}
	if t, err := time.Parse(time.RFC3339, str); err == nil {
		return t
	}

	return time.Time{}
}
