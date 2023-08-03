package storage

import (
	"devious/internal/config"
	"encoding/gob"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

func Get(localPath string, conf config.Config, gitDir string) error {
	/// Get metadata of specified file
	metadataFile, err := os.Open(localPath + MetaFileExtension)
	if err != nil {
		slog.Error("No metadata for file", slog.String("path", localPath))
		return err
	}

	// Decode metadata
	var metadata Metadata
	err = gob.NewDecoder(metadataFile).Decode(&metadata)
	if err != nil {
		slog.Error("Failed to decode metadata", slog.String("path", localPath))
		return err
	}

	// Get storage path
	storagePath := filepath.Join(conf.StorageDir, metadata.FileHash) + StorageFileExtension

	// Copy file to destination
	Copy(storagePath, localPath, conf)

	return nil
}
