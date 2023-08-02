package config

import (
	"devious/internal/git"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

type Config struct {
	StorageDir string `yaml:"storage-dir"`
}

var ConfigFileName = ".devious"

func Load() (Config, error) {
	// Create a new file at git root
	gitDir, err := git.GetRootDir()
	if err != nil {
		return Config{}, err
	}

	// Read the config file
	configFileContents, err := os.ReadFile(filepath.Join(gitDir, ConfigFileName))
	if err != nil {
		return Config{}, err
	}

	// Decode the config file
	var config Config
	err = yaml.Unmarshal(configFileContents, &config)
	if err != nil {
		return Config{}, err
	}

	slog.Info("Loaded config", slog.String("storage-dir", config.StorageDir))

	return config, nil
}
