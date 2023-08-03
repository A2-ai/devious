package cmd

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

func runInitCmd(cmd *cobra.Command, args []string) error {
	storageDir := args[0]

	// Convert storage dir to absolute path
	storageDir, err := filepath.Abs(storageDir)
	if err != nil {
		slog.Error("Failed to convert to absolute path", slog.String("path", storageDir))
		return err
	}

	// Get repository root
	gitDir, err := git.GetRootDir()
	if err != nil {
		return err
	}

	// Create a new file at git root
	dvsFile, err := os.Create(filepath.Join(gitDir, config.ConfigFileName))
	if err != nil {
		return err
	}
	defer dvsFile.Close()

	// Create config
	config := config.Config{
		StorageDir: storageDir,
	}

	// Convert config to YAML
	configYaml, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Write the default config to the file as YAML
	_, err = dvsFile.Write([]byte(configYaml))
	if err != nil {
		return err
	}

	slog.Info("Initialized devious", slog.String("storage-dir", storageDir))

	return nil
}

func getInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <storage-dir>",
		Short: "Initialize devious in the current git repository with the provided storage directory",
		Args:  cobra.ExactArgs(1),
		RunE:  runInitCmd,
	}

	return cmd
}
