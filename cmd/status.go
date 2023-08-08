package cmd

import (
	"dvs/internal/git"
	"dvs/internal/meta"
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func rawLog(args ...any) {
	os.Stdout.Write([]byte(fmt.Sprintln(args...)))
}

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

	for _, path := range metaPaths {
		relPath, err := git.GetRelativePath(".", path)
		if err != nil {
			return err
		}

		// Get file info
		metadata, err := meta.LoadFile(path)
		if err != nil {
			slog.Warn("File not in devious", slog.String("path", path))
			return err
		}

		cGray := color.New(color.Faint)

		// Check if file is available locally
		fileHash := fmt.Sprintf("%x", blake3.Sum256([]byte(localPath)))
		if metadata.FileHash

		// Determine file tag based on file status
		var fileTag string
		if metadata {
			fileTag = cGray.Sprint("deleted")
		} else if metadata.IsModified {
			fileTag = cGray.Sprint("modified")
		} else {
			fileTag = cGray.Sprint("unchanged")
		}

		rawLog("   ", fileTag, cGray.Sprint(humanize.Bytes(metadata.FileSize)))
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
