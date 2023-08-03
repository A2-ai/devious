package config

import (
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

type Config struct {
	StorageDir string `yaml:"storage-dir"`
}

var ConfigFileName = ".dvs.yaml"

func Load(gitDir string) (Config, error) {
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

	slog.Debug("Loaded config", slog.String("storage-dir", config.StorageDir))

	return config, nil
}
