package cmd

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/storage"
	"fmt"

	"github.com/spf13/cobra"
)

func runRemoveCmd(cmd *cobra.Command, args []string) error {
	// Get git dir
	gitDir, err := git.GetNearestRepoDir(".")
	if err != nil {
		return err
	}

	// Load the conf
	conf, err := config.Read(gitDir)
	if err != nil {
		log.PrintNotInitialized()
		log.DumpAndExit(0)
	}

	// Get flags
	dry, err := cmd.Flags().GetBool("dry")
	if err != nil {
		return err
	}

	// Remove each path from storage
	for i, path := range args {
		log.Print(fmt.Sprint(i+1)+"/"+fmt.Sprint(len(args)), " ", log.ColorFile(path))

		err = storage.Remove(path, conf, gitDir, dry)
		if err != nil {
			log.Print(log.ColorRed("    ✗"), "Failed to remove file", log.ColorFaint(err.Error()))
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "error",
				Message:  "failed to remove from storage",
				Location: path,
			})
		}
	}

	return nil
}

func getRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <file> <another-file> ...",
		Short: "Removes file(s) from storage and devious",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		RunE: runRemoveCmd,
	}

	cmd.Flags().BoolP("dry", "d", false, "run without actually removing files")

	return cmd
}
