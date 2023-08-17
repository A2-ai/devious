package meta

import (
	"dvs/internal/log"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var FileExtension = ".dvsmeta"

// Gets a list of all meta file paths in the directory recursively
func GetAllMetaFiles(dir string) (metaFiles []string, err error) {
	metaFiles = make([]string, 0)
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a meta file
		if filepath.Ext(path) == FileExtension {
			// strip extension
			metaFiles = append(metaFiles, path[0:len(path)-len(FileExtension)])
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return metaFiles, nil
}

// Parse globs into meta file paths, ignoring duplicates
func ParseGlobs(globs []string) (metaFiles []string) {
	metaFiles = make([]string, 0)

	for _, glob := range globs {
		// Remove meta file extension
		glob = strings.ReplaceAll(glob, FileExtension, "")

		// Skip if already queued
		if slices.Contains(metaFiles, glob) {
			continue
		}

		// Skip directories
		if s, err := os.Stat(glob); err == nil && s.IsDir() {
			log.Print(log.ColorBold(log.ColorYellow("!")), "Skipping directory", log.ColorFile(glob), "\n")
			continue
		}

		// Add to queue
		metaFiles = append(metaFiles, glob)
	}

	return metaFiles
}
