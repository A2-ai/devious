package cmd

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/meta"
	"dvs/internal/storage"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func runGetCmd(cmd *cobra.Command, args []string) error {
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
	recurse, err := cmd.Flags().GetBool("recurse")
	if err != nil {
		return err
	}

	dry, err := cmd.Flags().GetBool("dry")
	if err != nil {
		return err
	}

	var filesToGet []string

	// Parse each path
	for _, path := range args {
		// If the path is a directory, get all files in the directory
		if pathInfo, err := os.Stat(path); err == nil && pathInfo.IsDir() {
			var metaFiles []string
			if recurse {
				// Get all files in the current directory and subdirectories
				metaFiles, err = meta.GetAllMetaFiles(path)
				if err != nil {
					return err
				}
			} else {
				// Get all files in the current directory
				metaFiles, err = meta.GetMetaFiles(path)
				if err != nil {
					return err
				}
				if len(metaFiles) == 0 {
					absPath, _ := filepath.Abs(path)
					log.RawLog(log.ColorYellow("⚠"), "No devious files found in directory, skipping", log.ColorFile(absPath), "\n")
				}
			}
			filesToGet = append(filesToGet, metaFiles...)
		} else {
			filesToGet = append(filesToGet, path)
		}
	}

	// Get the queued files
	for i, file := range filesToGet {
		log.RawLog(fmt.Sprint(i+1)+"/"+fmt.Sprint(len(filesToGet)), " ", log.ColorFile(file))

		err = storage.Get(file, conf.StorageDir, gitDir, dry)
		if err != nil {
			log.RawLog(log.ColorRed("    ✘"), "Failed to get file", log.ColorFaint(err.Error()), "\n")
		} else {
			log.OverwritePreviousLine()
			log.RawLog("    Cleaning up...", log.ColorGreen("✔\n"))
		}
	}

	return nil
}

func getGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <file> <another-file> ...",
		Short: "Gets file(s) from storage",
		Long:  `Gets file(s) from storage. If no arguments are provided, get all files in the current directory.`,
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		RunE: runGetCmd,
	}

	cmd.Flags().BoolP("recurse", "r", false, "include subdirectories")
	cmd.Flags().BoolP("dry", "d", false, "run without actually getting files")

	return cmd
}
