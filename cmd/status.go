package cmd

import (
	"devious/internal/git"
	"devious/internal/meta"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func runStatusCmd(cmd *cobra.Command, args []string) error {
	// If no arguments are provided, get the status of all files in the current git repository
	if len(args) == 0 {
		// Get git dir
		gitDir, err := git.GetRootDir()
		if err != nil {
			return err
		}

		metaPaths, err := meta.GetMetaFiles(gitDir)
		if err != nil {
			return err
		}

		for _, metaPath := range metaPaths {
			metadata, err := meta.LoadFile(metaPath)
			if err != nil {
				return err
			}

			slog.Info("File status", slog.String("path", metaPath), slog.String("hash", metadata.FileHash))
		}
	}

	return nil
}

func getStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [file]",
		Short: "Gets the status of devious files in the current git repository, or a specific file if specified",
		RunE:  runStatusCmd,
	}

	return cmd
}
