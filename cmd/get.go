package cmd

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/meta"
	"dvs/internal/storage"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func runGetCmd(cmd *cobra.Command, args []string) error {
	defer log.Dump(log.JsonLogger)

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

	// Parse each glob
	var queuedPaths []string
	for _, glob := range args {
		// Strip meta file extension
		glob = strings.TrimSuffix(glob, meta.FileExtension)

		// Skip if already queued
		if slices.Contains(queuedPaths, glob) {
			continue
		}

		// Skip directories
		if s, err := os.Stat(glob); err == nil && s.IsDir() {
			continue
		}

		// Add to queue
		queuedPaths = append(queuedPaths, glob)
	}

	// Get the queued files
	for i, path := range queuedPaths {
		log.Print(fmt.Sprint(i+1)+"/"+fmt.Sprint(len(queuedPaths)), " ", log.ColorFile(path))

		err = storage.Get(path, conf.StorageDir, gitDir, dry)
		if err != nil {
			log.Print(log.ColorRed("    ✘"), "Failed to get file", log.ColorFaint(err.Error()), "\n")
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "error",
				Message:  "failed to get file",
				Location: path,
			})
		}
	}

	// Warn if no files were queued
	if len(queuedPaths) == 0 {
		log.Print(log.ColorBold(log.ColorYellow("!")), "No files were queued")
		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "warning",
			Message:  "no files were queued",
			Location: ".",
		})
	}

	return nil
}

func getGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <glob> [another-glob] ...",
		Short: "Gets file(s) from storage",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		RunE: runGetCmd,
	}

	cmd.Flags().BoolP("dry", "d", false, "run without actually getting files")

	return cmd
}
