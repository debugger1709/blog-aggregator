package main

import (
	"context"
	"fmt"
	"time"

	"github.com/debugger1709/blog-aggregator/internal/api"
	"github.com/debugger1709/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregate(s *state, cmd command) error {
	feed, err := api.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("a feed name and a url is required")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
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
	return nil
}
