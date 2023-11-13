package migrate

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/storage"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func migrateStorageFile(storageDir string, path string, dry bool) (match bool, err error) {
	// Ensure the file has the correct extension
	ext := filepath.Ext(path)
	if ext != "" {
		if dry {
			return true, nil
		}

		newPath := strings.TrimSuffix(path, ext)
		err := os.Rename(path, newPath)
		if err != nil {
			return false, err
		}

		path = newPath
		match = true
	}

	// Ensure the file is not at the root level
	if filepath.Dir(path) == storageDir {
		if dry {
			return true, nil
		}

		fileHash := filepath.Base(path)

		// Create new directory for first segment of the filename
		newDir := filepath.Join(storageDir, fileHash[:2])
		err := os.MkdirAll(newDir, storage.StorageDirPermissions)
		if err != nil {
			return false, err
		}

		// Move to correct location
		newPath := storage.GetStoragePath(storageDir, filepath.Base(path))
		err = os.Rename(path, newPath)

		match = true

		if os.IsExist(err) {
			// File already exists, delete the old one
			err = os.Remove(path)
			if err != nil {
				return match, err
			}
		} else if err != nil {
			return match, err
		}
	}

	return match, nil
}

// Migrates storage files to the latest format, returning a list of files that were modified
func MigrateStorageFiles(dry bool) (files []string, err error) {
	repoDir, _ := git.GetNearestRepoDir(".")
	config, err := config.Read(repoDir)
	if err != nil {
		return nil, err
	}

	// Iterate over all files in the storage directory
	err = filepath.WalkDir(config.StorageDir, func(path string, d fs.DirEntry, _ error) error {
		// Don't migrate the directories
		if d.IsDir() {
			return nil
		}

		// TODO respect gitignore?
		modified, err := migrateStorageFile(config.StorageDir, path, dry)
		if err != nil {
			return err
		}

		if modified {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
