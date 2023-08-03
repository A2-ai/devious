package cmd

import (
	"devious/internal/config"
	"devious/internal/git"
	"devious/internal/meta"
	"devious/internal/storage"
	"os"

	"github.com/spf13/cobra"
)

func runGetCmd(cmd *cobra.Command, args []string) error {
	// Get git dir
	gitDir, err := git.GetRootDir()
	if err != nil {
		return err
	}

	// Load the conf
	conf, err := config.Load(gitDir)
	if err != nil {
		return err
	}

	var filesToGet []string

	// If no arguments are provided, get all files in the current directory
	if len(args) == 0 {
		// Get current directory
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		recurse, err := cmd.Flags().GetBool("recurse")
		if err != nil {
			return err
		}

		if recurse {
			// Get all files in the current directory and subdirectories
			filesToGet, err = meta.GetAllMetaFiles(wd)
			if err != nil {
				return err
			}
		} else {
			// Get all files in the current directory
			filesToGet, err = meta.GetMetaFiles(wd)
			if err != nil {
				return err
			}
		}
	} else {
		filesToGet = args
	}

	// Get each file from storage
	for _, path := range filesToGet {
		storage.Get(path, conf, gitDir)
	}

	return nil
}

func getGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a file from storage",
		RunE:  runGetCmd,
	}

	cmd.Flags().BoolP("recurse", "r", false, "include subdirectories")

	return cmd
}
