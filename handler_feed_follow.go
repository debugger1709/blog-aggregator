package main

import (
	"context"
	"fmt"
	"time"

	"github.com/debugger1709/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("a url is required")
	}
	feedURL := cmd.args[0]
	feed, err := s.db.GetFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("cannot find feed from url %s %w", feedURL, err)
	}
	return createFeedFollow(user, feed, s)
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("issue when getting user %s from database %w", s.cfg.CurrentUserName, err)
		}
		return handler(s, cmd, user)
	}
}

func createFeedFollow(user database.User, feed database.Feed, s *state) error {
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	f, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("cannot create feed follow for user %s with url %s %w", s.cfg.CurrentUserName, feed.Url, err)
	}
	fmt.Println(f)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	rows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("cannot get all feeds for user %s %w", user.Name, err)
	}
	for _, row := range rows {
		fmt.Println(row.FeedName)
	}
	return nil
}

func handlerDeleteFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("an url is required")
	}
	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    cmd.args[0],
	}
	err := s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("issue when deleting feed follow for user %s with url %s %w", user.Name, cmd.args[0], err)
	}
	return nil
}
