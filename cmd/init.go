package cmd

import (
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/storage"

	"github.com/spf13/cobra"
)

func getInitRunner(cmd *cobra.Command, args []string) error {
	defer log.Dump(log.JsonLogger)

	// Get git root
	gitDir, err := git.GetNearestRepoDir(".")
	if err != nil {
		log.Print(log.ColorRed("✘"), "Failed to get git root", log.ColorFaint(err.Error()))
		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "error",
			Message:  "failed to get git root",
		})
		log.DumpAndExit(1)
	}

	// Initialize
	err = storage.Init(gitDir, args[0])
	if err != nil {
		log.Print(log.ColorRed("✘"), "Failed to initialize devious", log.ColorFaint(err.Error()))
		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "error",
			Message:  "failed to initialize devious",
		})
		log.DumpAndExit(1)
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
			log.PrintLogo()
		},
		RunE: getInitRunner,
	}

	return cmd
}
