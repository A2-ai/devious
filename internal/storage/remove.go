package storage

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/meta"
	"os"
	"path/filepath"
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
		log.Print("    Removing file from storage...")
		err = os.Remove(storagePath)
		if err == nil {
			log.OverwritePreviousLine()
			log.Print("    Removing file from storage...", log.ColorGreen("✔"))
		} else {
			log.OverwritePreviousLine()
			log.Print("    Removing file from storage...", log.ColorRed("✘ not present"))
		}
	} else {
		log.Print("    Removing file from storage...", log.ColorFaint("skipped (dry run)"))
	}

	// Remove path from gitignore
	err = git.RemoveIgnoreEntry(gitDir, path, dry)
	if err == nil {
		log.Print("    Removing gitignore entry...", log.ColorGreen("✔"))
	} else {
		log.Print("    Removing gitignore entry...", log.ColorRed("✘"))
	}

	// Remove metadata file
	if !dry {
		log.Print("    Removing metadata file...")
		err = os.Remove(path + meta.FileExtension)
		if err != nil {
			return err
		}
		log.OverwritePreviousLine()
		log.Print("    Removing metadata file...", log.ColorGreen("✔\n"))
	}

	return nil
}
