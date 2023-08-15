package storage

import (
	"dvs/internal/config"
	"dvs/internal/log"
	"os"
	"path/filepath"
)

var defaultDirPermissions = os.FileMode(0766)

func Init(rootDir string, storageDir string) error {
	// Get storage directory as absolute path
	storageDir, err := filepath.Abs(storageDir)
	if err != nil {
		log.Print(log.ColorRed("✘"), "Failed to convert destination to absolute path", log.ColorFile(storageDir))
		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "error",
			Message:  "failed to convert destination to absolute path",
			Location: storageDir,
		})

		return err
	}

	// Check storage dir permissions, and create if it doesn't exist
	fileInfo, err := os.Stat(storageDir)
	if err != nil {
		// Create storage dir and necessary parents
		err = os.MkdirAll(storageDir, defaultDirPermissions)
		if err != nil {
			log.Print(log.ColorRed("✘"), "Failed to create storage directory", log.ColorFile(storageDir))
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "error",
				Message:  "failed to create storage directory",
				Location: storageDir,
			})

			return err
		}

		log.Print(log.ColorGreen("✔"), "Created storage directory")
		log.JsonLogger.Actions = append(log.JsonLogger.Actions, log.JsonAction{
			Action: "created storage directory",
			Path:   storageDir,
		})
	} else {
		// Ensure destination is a directory
		if !fileInfo.IsDir() {
			log.Print(log.ColorRed("✘"), "Destination isn't a directory", log.ColorFile(storageDir))
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "error",
				Message:  "destination isn't a directory",
				Location: storageDir,
			})

			return err
		}

		// Warn if destination is not empty
		dir, err := os.ReadDir(storageDir)
		if err != nil {
			log.Print(log.ColorRed("✘"), "Failed to read storage directory", log.ColorFile(storageDir))
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "error",
				Message:  "failed to read storage directory",
				Location: storageDir,
			})
		} else if len(dir) > 0 {
			log.Print(log.ColorBold(log.ColorYellow("!")), "Storage directory isn't empty\n")
			log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
				Severity: "warning",
				Message:  "storage directory not empty",
				Location: storageDir,
			})
		}
	}

	// Write config
	err = config.Write(config.Config{
		StorageDir: storageDir,
	}, rootDir)
	if err != nil {
		return err
	}

	log.Print(log.ColorGreen("✔"), "Wrote config", log.ColorFile(filepath.Join(rootDir, config.ConfigFileName)))
	log.Print("    storage dir", log.ColorFile(storageDir))
	log.JsonLogger.Actions = append(log.JsonLogger.Actions, log.JsonAction{
		Action: "wrote config",
		Path:   filepath.Join(rootDir, config.ConfigFileName),
	})

	return nil
}
