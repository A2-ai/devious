package cmd

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/storage"

	"github.com/spf13/cobra"
)

func runRemoveCmd(cmd *cobra.Command, args []string) error {
	// Get git dir
	gitDir, err := git.GetRootDir()
	if err != nil {
		return err
	}

	// Load the conf
	conf, err := config.Read(gitDir)
	if err != nil {
		log.ThrowNotInitialized()
	}

	// Get flags
	dry, err := cmd.Flags().GetBool("dry")
	if err != nil {
		return err
	}

	// Remove each path from storage
	for _, path := range args {
		storage.Remove(path, conf, gitDir, dry)
	}

	return nil
}

func getRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <file> <another-file> ...",
		Short: "Removes file(s) from storage and devious",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runRemoveCmd,
	}

	cmd.Flags().BoolP("dry", "d", false, "run without actually removing files")

	return cmd
}
