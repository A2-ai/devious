package cmd

import (
	"dvs/internal/git"
	"dvs/internal/storage"
	"os"

	"github.com/spf13/cobra"
)

func runInitCmd(cmd *cobra.Command, args []string) error {
	// Get git root
	gitDir, err := git.GetNearestRepoDir(".")
	if err != nil {
		os.Exit(1)
		return err
	}

	// Initialize
	err = storage.Init(gitDir, args[0])
	if err != nil {
		os.Exit(1)
		return err
	}

	return nil
}

func getInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <storage-dir>",
		Short: "Initializes devious in the current git repository with the provided storage directory",
		Args:  cobra.ExactArgs(1),
		RunE:  runInitCmd,
	}

	return cmd
}
