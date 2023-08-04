package cmd

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/storage"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func runAddCmd(cmd *cobra.Command, args []string) error {
	// Get git dir
	gitDir, err := git.GetRootDir()
	if err != nil {
		return err
	}

	// Load the conf
	conf, err := config.Load(gitDir)
	if err != nil {
		log.ThrowNotInitialized()
	}

	// Get flags
	dry, err := cmd.Flags().GetBool("dry")
	if err != nil {
		return err
	}

	// Add each path to storage
	for _, path := range args {
		// if the path is a glob, get all files that match the glob
		// otherwise, add the file
		if strings.Contains(path, "*") {
			files, err := filepath.Glob(path)
			if err != nil {
				return err
			}

			for _, file := range files {
				storage.Add(file, conf, gitDir, dry)
			}

			continue
		} else {
			storage.Add(path, conf, gitDir, dry)
		}
	}

	return nil
}

func getAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <file> <another-file> <glob-pattern> ...",
		Short: "Add file(s) to storage",
		RunE:  runAddCmd,
	}

	cmd.Flags().BoolP("recurse", "r", false, "include subdirectories")
	cmd.Flags().BoolP("dry", "d", false, "run without actually adding files")

	return cmd
}
