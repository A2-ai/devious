package cmd

import (
	"devious/internal/git"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	StorageDir string `yaml:"storage-dir"`
}

func runInitCmd(cmd *cobra.Command, args []string) error {
	// Get the desired storage directory from the command line
	storageDir, err := cmd.Flags().GetString("storage-dir")
	if err != nil {
		return err
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create a new file at git root
	gitDir, err := git.GetRepositoryRoot(cwd)
	if err != nil {
		return err
	}

	// Create a new file at git root
	dvsFile, err := os.Create(gitDir + "/.devious")
	if err != nil {
		return err
	}
	defer dvsFile.Close()

	// Create config
	config := Config{
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
