package cmd

import (
	"devious/internal/config"
	"devious/internal/storage"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func runGetCmd(cmd *cobra.Command, args []string) error {
	localPath := args[0]

	// Load the conf
	conf, err := config.Load()
	if err != nil {
		return err
	}

	// Get metadata of desired file
	metadataFile, err := os.Open(args[0] + storage.MetaFileExtension)
	if err != nil {
		return err
	}

	// Decode metadata
	var metadata storage.Metadata
	err = gob.NewDecoder(metadataFile).Decode(&metadata)
	if err != nil {
		return err
	}

	// Copy file to destination
	// Open source file
	storagePath := filepath.Join(conf.StorageDir, metadata.FileHash) + storage.StorageFileExtension
	src, err := os.Open(storagePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(localPath)
	if err != nil {
		return err
	}

	defer dst.Close()

	// Copy the file to the storage directory
	copiedBytes, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	slog.Info("Copied file to local",
		slog.String("filesize", fmt.Sprintf("%1.2f MB", float64(copiedBytes)/1000000)),
		slog.String("from", storagePath), slog.String("to", localPath))

	return nil
}

func getGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a file from storage",
		RunE:  runGetCmd,
	}

	return cmd
}
