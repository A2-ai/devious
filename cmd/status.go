package cmd

import (
	"dvs/internal/git"
	"dvs/internal/meta"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func runStatusCmd(cmd *cobra.Command, args []string) error {
	var metaPaths []string

	// If no arguments are provided, get the status of all files in the current git repository
	if len(args) == 0 {
		// Get git dir
		gitDir, err := git.GetNearestRepoDir(".")
		if err != nil {
			return err
		}

		metaPaths, err = meta.GetAllMetaFiles(gitDir)
		if err != nil {
			return err
		}

		slog.Info("Total devious files", slog.Int("count", len(metaPaths)))
	} else {
		metaPaths = args
	}

	for _, metaPath := range metaPaths {
		metadata, err := meta.LoadFile(metaPath)
		if err != nil {
			slog.Warn("File not in devious", slog.String("path", metaPath))
			return err
		}

		slog.Info("File status", slog.String("path", metaPath), slog.String("hash", metadata.FileHash), slog.String("size", humanize.Bytes(metadata.FileSize)))
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
