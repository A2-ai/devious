package cmd

import (
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/storage"
	"os"

	"github.com/spf13/cobra"
)

func getInitRunner(cmd *cobra.Command, args []string) error {
	// Get git root
	gitDir, err := git.GetNearestRepoDir(".")
	if err != nil {
		os.Exit(1)
	}

	// Initialize
	err = storage.Init(gitDir, args[0])
	if err != nil {
		os.Exit(1)
	}

	return nil
}

func getInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <storage-dir>",
		Short: "Initializes devious in the current git repository with the provided storage directory",
		Long:  "Initializes devious in the current git repository with the provided storage directory. The storage directory should be a location accessible by all users of the repository.",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if json, err := cmd.Flags().GetBool("json"); err == nil && !json {
				log.PrintLogo()
			}
		},
		RunE: getInitRunner,
	}

	return cmd
}
