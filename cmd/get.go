package cmd

import (
	"devious/internal/config"
	"devious/internal/git"
	"devious/internal/storage"

	"github.com/spf13/cobra"
)

func runGetCmd(cmd *cobra.Command, args []string) error {
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

	// Get each file from storage
	for _, localPath := range args {
		storage.Get(localPath, conf, gitDir)
	}

	return nil
}

func getGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a file from storage",
		RunE:  runGetCmd,
	}

	return cmd
}
