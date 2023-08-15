package cmd

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/storage"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func runAddCmd(cmd *cobra.Command, args []string) error {
	defer log.Dump(log.JsonLogger)

	// Get flags
	dry, err := cmd.Flags().GetBool("dry")
	if err != nil {
		return err
	}

	message, err := cmd.Flags().GetString("message")
	if err != nil {
		return err
	}

	// Get git dir
	gitDir, err := git.GetNearestRepoDir(".")
	if err != nil {
		log.Print(log.ColorRed("✘"), "Couldn't find a parent git repository")
		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "error",
			Message:  "couldn't find a parent git repository",
			Location: ".",
		})
		log.DumpAndExit(0)
	}

	// Load the conf
	conf, err := config.Read(gitDir)
	if err != nil {
		log.PrintNotInitialized()
		log.DumpAndExit(0)
	}

	// Queue file paths
	var filesToAdd []string
	for _, path := range args {
		globMatches, err := filepath.Glob(path)

		if err != nil {
			log.Print(log.ColorBold(log.ColorYellow("!")), "Skipping invalid path", log.ColorFile(path), "\n")
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "warning",
				Message:  fmt.Sprintf("Skipping invalid path %s", path),
				Location: path,
			})
			continue
		}

		if !slices.Contains(globMatches, path) {
			log.Print(log.ColorBold(log.ColorYellow("!")), "Skipping invalid path", log.ColorFile(path), "\n")
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "warning",
				Message:  fmt.Sprintf("Skipping invalid path %s", path),
				Location: path,
			})
		}

		for _, match := range globMatches {
			// Ensure file is inside of the git repo
			absPath, err := filepath.Abs(match)
			if err != nil {
				log.Print(log.ColorBold(log.ColorYellow("!")), "Skipping invalid path", log.ColorFile(match), "\n")
				log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
					Severity: "warning",
					Message:  "skipped invalid path",
					Location: match,
				})
				continue
			}
			if strings.TrimPrefix(absPath, gitDir) == absPath {
				log.Print(log.ColorBold(log.ColorYellow("!")), "Skipping file outside of git repository", log.ColorFile(match), "\n")
				log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
					Severity: "warning",
					Message:  "skipped file outside of git repository",
					Location: match,
				})
				continue
			}

			// Check if file exists
			fileStat, err := os.Stat(match)
			if err != nil {
				log.Print(log.ColorBold(log.ColorYellow("!")), "Skipping invalid path", log.ColorFile(match), "\n")
				log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
					Severity: "warning",
					Message:  "skipped invalid path",
					Location: match,
				})

				continue
			}

			// Skip directories
			if fileStat.IsDir() {
				log.Print(log.ColorBold(log.ColorYellow("!")), "Skipping directory", log.ColorFile(match), "\n")
				log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
					Severity: "warning",
					Message:  "skipped directory",
					Location: match,
				})

				continue
			}

			// Add the file to the queue
			filesToAdd = append(filesToAdd, match)
		}
	}

	// Add each file to storage
	for i, file := range filesToAdd {
		log.Print(fmt.Sprint(i+1)+"/"+fmt.Sprint(len(filesToAdd)), " ", log.ColorFile(file))

		_, err := storage.Add(file, conf.StorageDir, gitDir, message, dry)
		if err != nil {
			log.Print(log.ColorRed("    ✗"), "Failed to add file", log.ColorFaint(err.Error()))
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "error",
				Message:  "error adding to storage",
				Location: file,
			})
			return err
		}
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
	cmd.Flags().StringP("message", "m", "", "add a message to the file's metadata")

	return cmd
}
