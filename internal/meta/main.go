package meta

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

var FileExtension = ".dvsmeta"

type Metadata struct {
	FileHash string `yaml:"file-hash"`
	FileSize uint64 `yaml:"file-size"`
}

// Creates a metadata file
func CreateFile(metadata Metadata, filePath string, dry bool) (err error) {
	var metadataFile *os.File
	if dry {
		slog.Debug("Dry run: creating metadata file", slog.String("path", filePath))
		return nil
	}

	slog.Debug("Creating metadata file", slog.String("path", filePath))

	metadataFile, err = os.Create(filePath + FileExtension)
	if err != nil {
		slog.Error("Failed to create metadata file", slog.String("path", filePath))
		return err
	}
	defer metadataFile.Close()

	err = json.NewEncoder(metadataFile).Encode(metadata)
	if err != nil {
		slog.Error("Failed to encode metadata", slog.String("path", filePath))
		return err
	}

	return nil
}

// Loads a metadata file
func LoadFile(path string) (metadata Metadata, err error) {
	metadataFile, err := os.Open(path + FileExtension)
	if err != nil {
		return metadata, err
	}

	err = json.NewDecoder(metadataFile).Decode(&metadata)
	if err != nil {
		slog.Error("Failed to decode metadata", slog.String("path", path))
		return metadata, err
	}

	return metadata, nil
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
