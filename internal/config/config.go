package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func getFileName() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("filed to get user home dir")
		return ""
	}

	filename := home + "/.gatorconfig.json"
	return filename
}

func Read() Config {
	filename := getFileName()
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("failed to read file: %v\n", filename)
		return Config{}
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("failed to unmarshal data from file")
		return Config{}
	}

	return config
}

func (c *Config) SetUser(username string) {
	c.Current_user_name = username

	json_data, err := json.Marshal(c)
	if err != nil {
		fmt.Println("failed to marshal data")
		return
	}
	filename := getFileName()
	err = os.WriteFile(filename, json_data, 0644)
	if err != nil {
		fmt.Printf("failed to write file with err: %v", err)
		return
	}
}
