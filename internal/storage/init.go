package storage

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

var defaultDirPermissions = os.FileMode(0766)

func Init(storageDir string) error {
	// Get storage directory as absolute path
	storageDir, err := filepath.Abs(storageDir)
	if err != nil {
		slog.Error("Failed to convert destination to absolute path", slog.String("path", storageDir))
		return err
	}

	// Check storage dir permissions, and create if it doesn't exist
	fileInfo, err := os.Stat(storageDir)
	if err != nil {
		slog.Info("Creating storage directory", slog.String("path", storageDir), slog.String("permissions", defaultDirPermissions.String()))

		// Create storage dir
		err = os.MkdirAll(storageDir, defaultDirPermissions)
		if err != nil {
			slog.Error("Failed to create storage directory", slog.String("path", storageDir))
			return err
		}
		fileInfo, err = os.Stat(storageDir)
		if err != nil {
			return err
		}
	}

	// Ensure destination is a directory
	if !fileInfo.IsDir() {
		slog.Error("Destination is not a directory", slog.String("path", storageDir))
		return err
	}

	// Ensure storage dir has write permissions
	if fileInfo.Mode().Perm()&0200 == 0 {
		slog.Error("Destination does not have write permissions", slog.String("path", storageDir))
		return os.ErrPermission
	}

	// Warn if not empty
	dir, err := os.ReadDir(storageDir)
	if err != nil {
		slog.Error("Failed to read storage directory", slog.String("path", storageDir))
		return err
	}
	if len(dir) > 0 {
		slog.Warn("Storage directory isn't empty", slog.String("path", storageDir))
	}

	// Get repository root
	gitDir, err := git.GetRootDir()
	if err != nil {
		return err
	}

	// Write config
	err = config.Write(config.Config{
		StorageDir: storageDir,
	}, gitDir)
	if err != nil {
		return err
	}

	return nil
}
