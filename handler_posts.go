package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/debugger1709/blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		n, err := strconv.Atoi(cmd.args[0])
		if err == nil {
			limit = n
		}
	}
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("cannot get posts for user %s", user.Name)
	}
	fmt.Println(posts)
	return nil
}
