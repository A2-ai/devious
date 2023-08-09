package storage

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/meta"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

// Remove a file from storage
func Remove(path string, conf config.Config, gitDir string, dry bool) error {
	/// Get metadata of specified file
	metadata, err := meta.Load(path)
	if err != nil {
		slog.Error("No metadata found for file", slog.String("path", path))
		return err
	}

	// Get storage path
	storagePath := filepath.Join(conf.StorageDir, metadata.FileHash) + FileExtension

	// Remove file from storage
	if !dry {
		err = os.Remove(storagePath)
		if err != nil {
			slog.Error("Failed to remove file from storage", slog.String("path", storagePath))
			return err
		}
		slog.Info("Removed file from storage", slog.String("path", storagePath))
	} else {
		slog.Info("Dry run: removed file from storage", slog.String("path", storagePath))
	}

	// Remove path from gitignore
	err = git.RemoveIgnoreEntry(gitDir, path, dry)
	if err != nil {
		return err
	}

	// Remove metadata file
	if !dry {
		err = os.Remove(path + meta.FileExtension)
		if err != nil {
			slog.Error("Failed to remove metadata file", slog.String("path", path))
		}
		slog.Debug("Removed metadata file", slog.String("path", path))
	} else {
		slog.Debug("Dry run: removed metadata file", slog.String("path", path))
	}

	return nil
}
