package cli

import (
	"errors"
	"fmt"

	"github.com/mrjkey/aggregator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	name string
	args []string
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.args) < 1 {
		return errors.New("the login handler expects a single argument, the username.")
	}

	username := cmd.args[0]
	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to: %v", username)

	return nil
}
