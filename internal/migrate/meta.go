package migrate

import (
	"dvs/internal/git"
	"dvs/internal/meta"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Attempts to migrate a metafile format from version 1 to the latest, returns true if anything was migrated
func migrateMetaFormatV1(path string) (match bool, err error) {
	type LegacyMetadata struct {
		FileHash  string    `json:"blake3_checksum"`
		FileSize  uint64    `json:"file_size_bytes"`
		Timestamp time.Time `json:"timestamp"`
		Message   string    `json:"message"`
		User      string    `json:"user"`
	}

	metadataFile, err := os.Open(path)
	if err != nil {
		return false, err
	}

	var legacyMetadata LegacyMetadata
	err = json.NewDecoder(metadataFile).Decode(&legacyMetadata)
	if err != nil {
		return false, err
	}

	// Map the contents to the new format
	metadata := meta.Metadata{
		FileHash:  legacyMetadata.FileHash,
		FileSize:  legacyMetadata.FileSize,
		Timestamp: legacyMetadata.Timestamp,
		Message:   legacyMetadata.Message,
		SavedBy:   legacyMetadata.User,
	}

	// Save the new format
	meta.Save(metadata, path, false)

	return true, nil
}

// Migrate a single meta file to the latest format, returns true if anything was migrated
func migrateMetaFile(path string, dry bool) (match bool, errs error) {
	// Ensure the file has the correct extension
	ext := filepath.Ext(path)
	if ext != meta.FileExtension {
		if dry {
			return true, nil
		}

		newPath := strings.TrimSuffix(path, ext) + meta.FileExtension
		err := os.Rename(path, newPath)
		if err != nil {
			return false, err
		}

		path = newPath
		match = true
	}

	// Ensure the file has the correct content structure
	pathNoExt := strings.TrimSuffix(path, meta.FileExtension)
	_, err := meta.Load(pathNoExt)
	if err != nil {
		if dry {
			return true, nil
		}

		_, err = migrateMetaFormatV1(path)
		if err != nil {
			return false, err
		}

		match = true
	}

	return match, err
}

// Migrates local meta files to the latest format, returning a list of files that were modified
func MigrateMetaFiles(dry bool) (files []string, err error) {
	// Iterate over all files in the git repository
	repoDir, _ := git.GetNearestRepoDir(".")
	filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, _ error) error {
		// TODO respect gitignore?

		// Check if the file is a meta file of some format
		if filepath.Ext(path) == ".dvsmeta" || filepath.Ext(path) == ".dvs" {
			wasOtherFormat, err := migrateMetaFile(path, dry)
			if err != nil {
				return err
			}
			if wasOtherFormat {
				files = append(files, path)
			}
		}

		return nil
	})

	return files, nil
}
