package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mrjkey/aggregator/internal/config"
	"github.com/mrjkey/aggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	gator_config := config.Read()
	db, err := sql.Open("postgres", gator_config.Db_url)
	if err != nil {
		fmt.Println("error opening database")
		os.Exit(1)
	}

	dbQueries := database.New(db)

	s := state{
		db:  dbQueries,
		cfg: &gator_config,
	}

	comms := commands{
		mapping: make(map[string]func(*state, command) error),
	}

	comms.register("login", handlerLogin)
	comms.register("register", handlerRegister)
	comms.register("reset", handlerReset)

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

	err = comms.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// gator_config.SetUser("jared")
	// gator_config_2 := config.Read()
	// fmt.Println(gator_config_2)
}

type state struct {
	db  *database.Queries
	cfg *config.Config
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
		return errors.New("the login handler expects a single argument, the username")
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to: %v\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	username := cmd.args[0]
	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	// _, err := s.db.GetUser(context.Background(), username)

	user, err := s.db.CreateUser(context.Background(), args)
	if err != nil {
		return err
	}

	fmt.Printf("User has been created: %v\n", user.Name)

	s.cfg.SetUser(user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.RemoveAllUsers(context.Background())
	return err
}
