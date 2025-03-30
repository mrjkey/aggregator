package main

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mrjkey/aggregator/internal/config"
	"github.com/mrjkey/aggregator/internal/database"
)

func TestHandlerReset(t *testing.T) {
	// Mock state
	s := &state{
		db:  &mockDB{},
		cfg: &config.Config{},
	}

	err := handlerReset(s, command{})
	if err != nil {
		t.Errorf("handlerReset() error = %v", err)
	}
}

func TestHandlerUsers(t *testing.T) {
	// Mock state
	s := &state{
		db:  &mockDB{},
		cfg: &config.Config{Current_user_name: "testuser"},
	}

	err := handlerUsers(s, command{})
	if err != nil {
		t.Errorf("handlerUsers() error = %v", err)
	}
}

// Mock database implementation
type mockDB struct{}

func (m *mockDB) RemoveAllUsers(ctx context.Context) error {
	return nil
}

func (m *mockDB) RemoveAllFeeds(ctx context.Context) error {
	return nil
}

func (m *mockDB) GetUsers(ctx context.Context) ([]database.User, error) {
	return []database.User{
		{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "testuser",
		},
	}, nil
}
