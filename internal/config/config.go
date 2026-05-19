package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

func write(cfg *Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("cannot create file at %s: %w", configFilePath, err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return fmt.Errorf("cannot encode config file %s: %w", configFilePath, err)
	}
	return nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("cannot open file at %s: %w", configFilePath, err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("cannot decode config file %s: %w", configFilePath, err)
	}
	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get home directory %w", err)
	}
	fullPath := filepath.Join(homeDir, configFileName)
	return fullPath, nil
}
