package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "gatorconfig.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the temp file
	testData := `{"db_url": "test_db_url", "current_user_name": "test_user"}`
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Override the getFileName function to return the temp file name
	originalGetFileName := getFileName
	getFileName = func() string {
		return tempFile.Name()
	}
	defer func() { getFileName = originalGetFileName }()

	// Call the Read function
	config := Read()

	// Verify the results
	if config.Db_url != "test_db_url" {
		t.Errorf("Expected db_url to be 'test_db_url', got '%s'", config.Db_url)
	}
	if config.Current_user_name != "test_user" {
		t.Errorf("Expected current_user_name to be 'test_user', got '%s'", config.Current_user_name)
	}
}

func TestSetUser(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "gatorconfig.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write initial test data to the temp file
	initialData := `{"db_url": "test_db_url", "current_user_name": "initial_user"}`
	if _, err := tempFile.Write([]byte(initialData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Override the getFileName function to return the temp file name
	originalGetFileName := getFileName
	getFileName = func() string {
		return tempFile.Name()
	}
	defer func() { getFileName = originalGetFileName }()

	// Read the initial config
	config := Read()

	// Call the SetUser function
	newUser := "new_user"
	if err := config.SetUser(newUser); err != nil {
		t.Fatalf("Failed to set user: %v", err)
	}

	// Read the updated config
	updatedConfig := Read()

	// Verify the results
	if updatedConfig.Current_user_name != newUser {
		t.Errorf("Expected current_user_name to be '%s', got '%s'", newUser, updatedConfig.Current_user_name)
	}

	// Verify the file content
	fileContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	var fileConfig Config
	if err := json.Unmarshal(fileContent, &fileConfig); err != nil {
		t.Fatalf("Failed to unmarshal file content: %v", err)
	}
	if fileConfig.Current_user_name != newUser {
		t.Errorf("Expected file current_user_name to be '%s', got '%s'", newUser, fileConfig.Current_user_name)
	}
}
