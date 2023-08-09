package cmd

import (
	"dvs/internal/file"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/meta"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func runStatusCmd(cmd *cobra.Command, args []string) error {
	log.RawLog(color.New(color.Bold).Sprint("👺\n"))

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
	} else {
		metaPaths = args
	}

	// Create colors
	colorFaintBold := color.New(color.Faint, color.Bold)
	colorFilePulled := color.New(color.FgGreen, color.Bold)
	colorFileOutdated := color.New(color.FgHiYellow, color.Bold)
	colorFileNotPulled := color.New(color.FgRed, color.Bold)

	// Track number of files
	numFilesPulled := 0
	numFilesOutdated := 0
	numFilesNotPulled := 0

	// Print info about each file
	log.RawLog(color.New(color.Bold).Sprint("file info "),
		colorFilePulled.Sprint("●"), "up to date ",
		colorFileOutdated.Sprint("●"), "out of date ",
		colorFileNotPulled.Sprint("●"), "not present ",
	)
	for _, path := range metaPaths {
		relPath, err := git.GetRelativePath(".", path)
		if err != nil {
			return err
		}

		// Get file info
		metadata, err := meta.Load(path)
		if err != nil {
			slog.Warn("File not in devious", slog.String("path", path))
			return err
		}

		// Determine file tag based on file status
		var fileStatus string
		_, statErr := os.Stat(path)
		fileHash, hashErr := file.GetFileHash(path)
		if statErr == nil && hashErr == nil && fileHash == metadata.FileHash {
			fileStatus = colorFilePulled.Sprint("●")
			numFilesPulled++
		} else if statErr == nil {
			fileStatus = colorFileOutdated.Sprint("●")
			numFilesOutdated++
		} else {
			fileStatus = colorFileNotPulled.Sprint("●")
			numFilesNotPulled++
		}

		// Print file info
		log.RawLog("   ", fileStatus, colorFaintBold.Sprint(relPath), " ", color.New(color.Faint).Sprint(humanize.Bytes(metadata.FileSize)))
	}

	// Print overview
	log.RawLog(color.New(color.Bold).Sprint("\ntotals"))
	log.RawLog(colorFilePulled.Sprint(numFilesPulled), "up to date", colorFileOutdated.Sprint(numFilesOutdated), "out of date", colorFileNotPulled.Sprint(numFilesNotPulled), "not present", "\n")

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
