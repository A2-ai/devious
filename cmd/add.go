package cmd

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/storage"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func runAddCmd(cmd *cobra.Command, args []string) error {
	// Get git dir
	gitDir, err := git.GetNearestRepoDir(".")
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

	// Queue file paths
	var filesToAdd []string
	for _, path := range args {
		globMatches, err := filepath.Glob(path)
		// If the path is a glob pattern, add all matches
		// Otherwise, add the path itself
		if err == nil {
			for _, match := range globMatches {
				// Skip directories
				if fileStat, _ := os.Stat(match); fileStat.IsDir() {
					slog.Warn("Skipping directory", slog.String("path", match))
					continue
				}

				// Add the file to the queue
				filesToAdd = append(filesToAdd, match)
			}
		} else {
			filesToAdd = append(filesToAdd, path)
		}
	}

	// Add each file to storage
	for i, file := range filesToAdd {
		log.RawLog(fmt.Sprint(i+1)+"/"+fmt.Sprint(len(filesToAdd)), " ", log.ColorFile(file))
		storage.Add(file, conf.StorageDir, gitDir, dry)
	}

	return nil
}

func getAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <file> <another-file> <glob-pattern> ...",
		Short: "Adds file(s) to storage",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		RunE: runAddCmd,
	}

	cmd.Flags().BoolP("recurse", "r", false, "include subdirectories")
	cmd.Flags().BoolP("dry", "d", false, "run without actually adding files")

	return cmd
}
