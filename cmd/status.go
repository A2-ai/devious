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
		SavedBy   string `json:"savedBy"`
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
		// Get meta files from globs
		metaPaths = meta.ParseGlobs(args)
	}

	// Create colors
	colorFilePulled := color.New(color.FgGreen, color.Bold).Sprint
	colorFileOutdated := color.New(color.FgHiYellow, color.Bold).Sprint
	colorFileNotPulled := color.New(color.FgRed, color.Bold).Sprint
	iconPulled := colorFilePulled("✔")
	iconOutdated := colorFileOutdated("!")
	iconNotPulled := colorFileNotPulled("✘")

	// Track number of files
	numFilesPulled := 0
	numFilesOutdated := 0
	numFilesNotPulled := 0

	// Print info about each file
	if len(metaPaths) == 0 {
		log.Print(log.ColorBold(log.ColorYellow("!")), "No files were queued")
		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "warning",
			Message:  "no files queued",
			Location: ".",
		})
	} else {
		log.Print(color.New(color.Bold).Sprint("file info "),
			iconPulled, "up to date ",
			iconOutdated, "out of date ",
			iconNotPulled, "not present ",
		)
	}
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
			log.Print(log.ColorRed("    ✘"), "File not in devious", log.ColorFile(relPath))
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
			fileLight = iconPulled
			fileStatus = "up to date"
			numFilesPulled++
		} else if statErr == nil {
			fileLight = iconOutdated
			fileStatus = "out of date"
			numFilesOutdated++
		} else {
			fileLight = iconNotPulled
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
			log.ColorFaint(metadata.SavedBy), "",
			log.ColorFaint(metadata.Message),
		)
		jsonLogger.Files = append(jsonLogger.Files, jsonFileResult{
			Path:      relPath,
			Status:    fileStatus,
			FileSize:  metadata.FileSize,
			FileHash:  metadata.FileHash,
			Timestamp: timestamp,
			SavedBy:   metadata.SavedBy,
			Message:   metadata.Message,
		})
	}

	// Print overview
	if len(metaPaths) > 0 {
		log.Print(color.New(color.Bold).Sprint("\ntotals"))
		log.Print(
			colorFilePulled(numFilesPulled), "up to date ",
			colorFileOutdated(numFilesOutdated), "out of date ",
			colorFileNotPulled(numFilesNotPulled), "not present ",
		)
	}

	log.Dump(jsonLogger)

	return nil
}

func getStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [glob] [another-glob]",
		Short: "Gets the status of files tracked by devious",
		Long:  "Gets the status of files tracked by devious. If glob(s) are provided, only those globs will be checked. Otherwise, all files in the current git repository will be checked.",
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		RunE: runStatusCmd,
	}

	return cmd
}
