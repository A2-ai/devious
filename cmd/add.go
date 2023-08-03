package cmd

import (
	"devious/internal/config"
	"devious/internal/git"
	"devious/internal/storage"

	"github.com/spf13/cobra"
)

func runAddCmd(cmd *cobra.Command, args []string) error {
	// Get git dir
	gitDir, err := git.GetRootDir()
	if err != nil {
		return err
	}

	// Load the conf
	conf, err := config.Load(gitDir)
	if err != nil {
		return err
	}

	// For each file, add it to the storage
	for _, filePath := range args {
		storage.Add(filePath, conf, gitDir)
	}

	return nil
}

func getAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <file> <another-file> ...",
		Short: "Add file(s) to storage",
		RunE:  runAddCmd,
	}

	return cmd
}
