package config

import (
	"fmt"
	"github.com/AndreySirin/-Effective-Mobile-/internal/server"
	"github.com/AndreySirin/-Effective-Mobile-/internal/storage"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	Server   server.Config  `yaml:"server"`
	Postgres storage.Config `yaml:"postgres"`
}

func Load(lg *slog.Logger) (*Config, error) {
	lg.Info("loading config")

	exePath, err := os.Executable()
	if err != nil {
		lg.Error("failed to get executable path", "err", err)
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	configPath := filepath.Join(exeDir, "config.yaml")
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		lg.Error("failed to get absolute config path", "path", configPath, "err", err)
		return nil, fmt.Errorf("failed to get absolute config path: %w", err)
	}

	data, err := os.ReadFile(absConfigPath)
	if err != nil {
		lg.Error("failed to read config file", "path", absConfigPath, "err", err)
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		lg.Error("failed to parse config file", "path", absConfigPath, "err", err)
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	lg.Info("config loaded successfully", "path", absConfigPath)
	return &config, nil
}
