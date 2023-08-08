package meta

import (
	"io/fs"
	"os"
	"path/filepath"
)

var FileExtension = ".dvsmeta"

type Metadata struct {
	FileHash string `yaml:"file-hash"`
	FileSize uint64 `yaml:"file-size"`
}

// Gets a list of all meta file paths in the directory
func GetMetaFiles(dir string) (metaFiles []string, err error) {
	metaFiles = make([]string, 0)
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, dirEntry := range dirEntries {
		// Check if the file is a meta file
		if filepath.Ext(dirEntry.Name()) == FileExtension {
			metaFiles = append(metaFiles, filepath.Join(dir, dirEntry.Name()[0:len(dirEntry.Name())-len(FileExtension)]))
		}
	}

	return metaFiles, nil
}

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
