package cli

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mrjkey/aggregator/internal/config"
	"github.com/mrjkey/aggregator/internal/database"
)

type mockDB struct {
	users map[string]database.User
}

func (m *mockDB) GetUser(ctx context.Context, name string) (database.User, error) {
	user, exists := m.users[name]
	if !exists {
		return database.User{}, sql.ErrNoRows
	}
	return user, nil
}

func (m *mockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	user := database.User{
		ID:        arg.ID,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
		Name:      arg.Name,
	}
	m.users[arg.Name] = user
	return user, nil
}

func TestHandlerLogin(t *testing.T) {
	mockDB := &mockDB{
		users: map[string]database.User{
			"test_user": {
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      "test_user",
			},
		},
	}

	cfg := &config.Config{}
	s := &state{
		db:  mockDB,
		cfg: cfg,
	}

	cmd := command{
		name: "login",
		args: []string{"test_user"},
	}

	err := handlerLogin(s, cmd)
	if err != nil {
		t.Fatalf("handlerLogin failed: %v", err)
	}

	if s.cfg.Current_user_name != "test_user" {
		t.Errorf("Expected current_user_name to be 'test_user', got '%s'", s.cfg.Current_user_name)
	}
}

func TestHandlerRegister(t *testing.T) {
	mockDB := &mockDB{
		users: make(map[string]database.User),
	}

	cfg := &config.Config{}
	s := &state{
		db:  mockDB,
		cfg: cfg,
	}

	cmd := command{
		name: "register",
		args: []string{"new_user"},
	}

	err := handlerRegister(s, cmd)
	if err != nil {
		t.Fatalf("handlerRegister failed: %v", err)
	}

	if s.cfg.Current_user_name != "new_user" {
		t.Errorf("Expected current_user_name to be 'new_user', got '%s'", s.cfg.Current_user_name)
	}

	if _, exists := mockDB.users["new_user"]; !exists {
		t.Errorf("Expected user 'new_user' to be created in the database")
	}
}
