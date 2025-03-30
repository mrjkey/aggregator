package cli

import (
	"errors"
	"fmt"

	"github.com/mrjkey/aggregator/internal/config"
)

type state struct {
	Config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	mapping map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.mapping[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.mapping[cmd.name]
	if !ok {
		return errors.New("function not found")
	}

	return f(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
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
