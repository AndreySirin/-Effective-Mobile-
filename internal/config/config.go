package config

import (
	"fmt"
	"github.com/AndreySirin/-Effective-Mobile-/internal/server"
	"github.com/AndreySirin/-Effective-Mobile-/internal/storage"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Server   server.Config  `yaml:"server"`
	Postgres storage.Config `yaml:"postgres"`
}

func Load() (*Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	configPath := filepath.Join(exeDir, "config.example.yaml")

	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute config path: %w", err)
	}
	data, err := os.ReadFile(absConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	return &config, nil
}
