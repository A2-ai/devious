package cmd

import (
	"devious/internal/config"
	"devious/internal/git"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func runAddCmd(cmd *cobra.Command, args []string) error {
	// Get git dir
	gitDir, err := git.GetRootDir()
	if err != nil {
		return err
	}

	// Get desired file
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	// Create metadata file
	metadataFile, err := os.Create(filePath + config.ConfigFileName)
	if err != nil {
		return err
	}
	defer metadataFile.Close()

	// Write metadata to file
	_, err = metadataFile.Write([]byte("metadata lorem ipsum"))
	if err != nil {
		return err
	}

	// Load the config
	config, err := config.Load()
	if err != nil {
		return err
	}

	// Open source and destination files
	src, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer src.Close()

	dstPath := filepath.Join(config.StorageDir, filepath.Base(filePath))
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the file to the storage directory
	copiedBytes, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	slog.Info("Copied file to storage",
		slog.String("filesize", fmt.Sprintf("%1.2f MB", float64(copiedBytes)/1000000)),
		slog.String("from", filePath), slog.String("to", dstPath))

	// Add file to gitignore
	ignoreFile, err := os.OpenFile(filepath.Join(gitDir, ".gitignore"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer ignoreFile.Close()

	_, err = ignoreFile.WriteString("\n\n# Devious entry\n/" + filepath.Base(filePath))
	return err
}

func getAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a file to storage",
		RunE:  runAddCmd,
	}

	cmd.Flags().StringP("file", "f", "", "The file to add to storage")
	cmd.MarkFlagFilename("file")
	cmd.MarkFlagRequired("file")

	return cmd
}
