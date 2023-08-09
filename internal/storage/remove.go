package storage

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/log"
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
		return err
	}

	// Get storage path
	storagePath := filepath.Join(conf.StorageDir, metadata.FileHash) + FileExtension

	// Remove file from storage
	if !dry {
		log.RawLog("    Removing file from storage...")
		err = os.Remove(storagePath)
		if err == nil {
			log.OverwritePreviousLine()
			log.RawLog("    Removing file from storage...", log.ColorGreen("✔"))
		} else {
			log.RawLog("    Removing file from storage...", log.ColorRed("✘ not present"))
		}
	} else {
		slog.Info("Dry run: removed file from storage", slog.String("path", storagePath))
	}

	// Remove path from gitignore
	err = git.RemoveIgnoreEntry(gitDir, path, dry)
	if err == nil {
		log.RawLog("    Removing gitignore entry...", log.ColorGreen("✔"))
	} else {
		log.RawLog("    Removing gitignore entry...", log.ColorRed("✘"))
	}

	// Remove metadata file
	if !dry {
		log.RawLog("    Removing metadata file...")
		err = os.Remove(path + meta.FileExtension)
		if err != nil {
			return err
		}
		log.OverwritePreviousLine()
		log.RawLog("    Removing metadata file...", log.ColorGreen("✔\n"))
	} else {
		slog.Debug("Dry run: removed metadata file", slog.String("path", path))
	}

	return nil
}
