package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/debugger1709/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("a name is required")
	}
	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("this name does not exist in database %w", err)
		}
		return fmt.Errorf("error getting user in database %w", err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Printf("user %s has been set\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("a name is required")
	}
	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("there is already a user with this name")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking name in database %w", err)
	}
	userParam := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	}
	_, err = s.db.CreateUser(context.Background(), userParam)
	if err != nil {
		return fmt.Errorf("error creating user in database %w", err)
	}
	fmt.Printf("Success: User '%s' has been created successfully in database!\n", name)
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("user was created in database, but failed to log in via config: %w", err)
	}
	fmt.Printf("user %s has been set\n", name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("issue when deleting all users %w", err)
	}
	fmt.Printf("all users removed\n")
	return nil
}

func handlerAllUsers(s *state, cmd command) error {
	names, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("issue when getting all users %w", err)
	}
	current := s.cfg.CurrentUserName
	for _, name := range names {
		if name == current {
			fmt.Printf("* %s (current)\n", name)
			continue
		}
		fmt.Printf("* %s\n", name)
	}
	return nil
}
