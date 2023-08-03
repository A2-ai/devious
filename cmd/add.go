package cmd

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/storage"

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
	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}

	var filesToAdd []string

	// If all, add all files in the current directory
	if all {
		// recurse, err := cmd.Flags().GetBool("recurse")
		// if err != nil {
		// 	return err
		// }

		// // Get current directory
		// wd, err := os.Getwd()
		// if err != nil {
		// 	return err
		// }

		// if recurse {
		// 	// Add all files in the current directory and subdirectories
		// 	filesToAdd, err = meta.GetAllNonMetaFiles(wd)
		// 	if err != nil {
		// 		return err
		// 	}
		// } else {
		// 	// Add all files in the current directory
		// 	filesToAdd, err = meta.GetNonMetaFiles(wd)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	} else {
		filesToAdd = args
	}

	// Add each file to storage
	for _, path := range filesToAdd {
		storage.Add(path, conf, gitDir)
	}

	return nil
}

func getAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <file> <another-file> ...",
		Short: "Add file(s) to storage",
		RunE:  runAddCmd,
	}

	cmd.Flags().BoolP("all", "a", false, "include all files in the current directory")
	cmd.Flags().BoolP("recurse", "r", false, "include subdirectories")

	return cmd
}
