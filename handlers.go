package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("a username is required")
	}
	userName := cmd.args[0]
	err := s.cfg.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Printf("user %s has been set\n", userName)
	return nil
}
