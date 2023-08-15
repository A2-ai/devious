package cmd

import (
	"dvs/internal/file"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/meta"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func runStatusCmd(cmd *cobra.Command, args []string) error {
	type jsonFileResult struct {
		Path      string `json:"path"`
		Status    string `json:"status"`
		FileSize  uint64 `json:"fileSize"`
		FileHash  string `json:"fileHash"`
		Timestamp string `json:"timestamp"`
		User      string `json:"user"`
		Message   string `json:"message"`
	}

	type jsonResult struct {
		Files  []jsonFileResult `json:"files"`
		Errors []log.JsonIssue  `json:"errors"`
	}

	jsonLogger := jsonResult{}
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
		var fileLight string
		var fileStatus string
		_, statErr := os.Stat(path)
		fileHash, hashErr := file.GetFileHash(path)
		if statErr == nil && hashErr == nil && fileHash == metadata.FileHash {
			fileLight = colorFilePulled.Sprint("●")
			fileStatus = "up to date"
			numFilesPulled++
		} else if statErr == nil {
			fileLight = colorFileOutdated.Sprint("●")
			fileStatus = "out of date"
			numFilesOutdated++
		} else {
			fileLight = colorFileNotPulled.Sprint("●")
			fileStatus = "not present"
			numFilesNotPulled++
		}

		// Determine whether to print timestamp
		timestamp := ""
		if !metadata.Timestamp.IsZero() {
			timestamp = metadata.Timestamp.Format(time.DateTime)
		}

		// Print file info
		log.Print("   ", fileLight,
			log.ColorFile(relPath), "",
			log.ColorFaint(humanize.Bytes(metadata.FileSize)), "",
			log.ColorFaint(timestamp), "",
			log.ColorFaint(metadata.User), "",
			log.ColorFaint(metadata.Message),
		)
		jsonLogger.Files = append(jsonLogger.Files, jsonFileResult{
			Path:      relPath,
			Status:    fileStatus,
			FileSize:  metadata.FileSize,
			FileHash:  metadata.FileHash,
			Timestamp: timestamp,
			User:      metadata.User,
			Message:   metadata.Message,
		})
	}

	// Print overview
	log.Print(color.New(color.Bold).Sprint("\ntotals"))
	log.Print(
		colorFilePulled.Sprint(numFilesPulled), "up to date ",
		colorFileOutdated.Sprint(numFilesOutdated), "out of date ",
		colorFileNotPulled.Sprint(numFilesNotPulled), "not present ",
	)

	log.Dump(jsonLogger)

	return nil
}

func getStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [path]",
		Short: "Gets the status of files tracked by devious",
		Long:  "Gets the status of files tracked by devious. If path(s) are provided, only those files will be checked. Otherwise, all files in the current git repository will be checked.",
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		RunE: runStatusCmd,
	}

	return cmd
}
