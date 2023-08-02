package cmd

import (
	"devious/internal/config"
	"devious/internal/git"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func runInitCmd(cmd *cobra.Command, args []string) error {
	// Get the desired storage directory from the command line
	storageDir, err := cmd.Flags().GetString("storage-dir")
	if err != nil {
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

	return nil
}

func getInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize devious",
		RunE:  runInitCmd,
	}

	cmd.Flags().StringP("storage-dir", "d", "", "The directory where devious will store managed files")
	cmd.MarkFlagDirname("storage-dir")
	cmd.MarkFlagRequired("storage-dir")

	return cmd
}
