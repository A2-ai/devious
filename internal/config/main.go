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

func Read(rootDir string) (Config, error) {
	// Read the config file
	configFileContents, err := os.ReadFile(filepath.Join(rootDir, ConfigFileName))
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

func Write(config Config, dir string) error {
	// Create config file
	dvsFile, err := os.Create(filepath.Join(dir, ConfigFileName))
	if err != nil {
		slog.Error("Failed to create config file", slog.String("dir", dir))
		return err
	}
	defer dvsFile.Close()

	// Convert config to YAML
	configYaml, err := yaml.Marshal(config)
	if err != nil {
		slog.Error("Failed to convert config to YAML", slog.String("dir", dir))
		return err
	}

	// Write the default config to the file as YAML
	_, err = dvsFile.Write([]byte(configYaml))
	if err != nil {
		slog.Error("Failed to write config", slog.String("dir", dir))
		return err
	}

	return nil
}
