package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mrjkey/aggregator/internal/config"
)

func main() {
	gator_config := config.Read()
	s := state{
		Config: &gator_config,
	}

	comms := commands{
		mapping: make(map[string]func(*state, command) error),
	}

	comms.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Need more arguments!")
		os.Exit(1)
	}

	// fmt.Println(args)

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err := comms.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// gator_config.SetUser("jared")
	// gator_config_2 := config.Read()
	// fmt.Println(gator_config_2)
}

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

	fmt.Printf("User has been set to: %v\n", username)

	return nil
}
