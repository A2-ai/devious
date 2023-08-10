package storage

import (
	"dvs/internal/config"
	"dvs/internal/log"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

var defaultDirPermissions = os.FileMode(0766)

func Init(rootDir string, storageDir string) error {
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

		// Create storage dir and necessary parents
		err = os.MkdirAll(storageDir, defaultDirPermissions)
		if err != nil {
			slog.Error("Failed to create storage directory", slog.String("path", storageDir))
			return err
		}
	} else {
		// Ensure destination is a directory
		if !fileInfo.IsDir() {
			slog.Error("Destination isn't a directory", slog.String("path", storageDir))
			return err
		}

		// Warn if destination is not empty
		dir, err := os.ReadDir(storageDir)
		if err != nil {
			slog.Error("Failed to read storage directory", slog.String("path", storageDir))
		} else if len(dir) > 0 {
			log.RawLog(log.ColorYellow("⚠"), "Storage directory isn't empty\n")
		}
	}

	// Write config
	err = config.Write(config.Config{
		StorageDir: storageDir,
	}, rootDir)
	if err != nil {
		return err
	}

	log.RawLog(log.ColorGreen("✔"), "Created config", log.ColorFile(rootDir))
	log.RawLog(log.ColorGreen("✔"), "Initialized devious")
	log.RawLog("    storage dir", log.ColorFile(storageDir))

	return nil
}
