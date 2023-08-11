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
	} else {
		metaPaths = args
	}

	// Create colors
	colorFilePulled := color.New(color.FgGreen, color.Bold)
	colorFileOutdated := color.New(color.FgHiYellow, color.Bold)
	colorFileNotPulled := color.New(color.FgRed, color.Bold)

	// Track number of files
	numFilesPulled := 0
	numFilesOutdated := 0
	numFilesNotPulled := 0

	// Print info about each file
	log.Print(color.New(color.Bold).Sprint("file info "),
		colorFilePulled.Sprint("●"), "up to date ",
		colorFileOutdated.Sprint("●"), "out of date ",
		colorFileNotPulled.Sprint("●"), "not present ",
	)
	for _, path := range metaPaths {
		relPath, err := git.GetRelativePath(".", path)
		if err != nil {
			log.Print(log.ColorRed("\n✘"), "Failed to get relative path", log.ColorFaint(err.Error()))
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "error",
				Message:  "failed to get relative path",
				Location: path,
			})

			continue
		}

		// Get file info
		metadata, err := meta.Load(path)
		if err != nil {
			log.Print(log.ColorRed("\n✘"), "File not in devious", log.ColorFile(relPath))
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "error",
				Message:  "file not in devious",
				Location: relPath,
			})

			continue
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
		log.Print("   ", fileStatus, log.ColorFile(relPath), " ", log.ColorFaint(humanize.Bytes(metadata.FileSize)))
		log.JsonLogger.Files[relPath] = log.JsonFile{
			Action: "status",
			Status: fileStatus,
		}
	}

	// Print overview
	log.Print(color.New(color.Bold).Sprint("\ntotals"))
	log.Print(
		colorFilePulled.Sprint(numFilesPulled), "up to date ",
		colorFileOutdated.Sprint(numFilesOutdated), "out of date ",
		colorFileNotPulled.Sprint(numFilesNotPulled), "not present ",
	)

	return nil
}

func getStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [file]",
		Short: "Gets the status of files tracked by devious",
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		RunE: runStatusCmd,
	}

	return cmd
}
