package meta

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

var FileExtension = ".dvsmeta"

type Metadata struct {
	FileHash string `yaml:"file-hash"`
}

func CreateFile(metadata Metadata, filePath string) (err error) {
	metadataFile, err := os.Create(filePath + FileExtension)
	if err != nil {
		slog.Error("Failed to create metadata file", slog.String("path", filePath))
		return err
	}
	defer metadataFile.Close()

	err = gob.NewEncoder(metadataFile).Encode(metadata)
	if err != nil {
		slog.Error("Failed to encode metadata", slog.String("path", filePath))
		return err
	}

	return nil
}

func LoadFile(filePath string) (metadata Metadata, err error) {
	metadataFile, err := os.Open(filePath)
	if err != nil {
		slog.Error("No metadata for file", slog.String("path", filePath))
		return metadata, err
	}

	err = gob.NewDecoder(metadataFile).Decode(&metadata)
	if err != nil {
		slog.Error("Failed to decode metadata", slog.String("path", filePath))
		return metadata, err
	}

	return metadata, nil
}

// Gets a list of all meta file paths in the directory recursively
func GetMetaFiles(dir string) (metaFiles []string, err error) {
	metaFiles = make([]string, 0)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a meta file
		if filepath.Ext(path) == FileExtension {
			metaFiles = append(metaFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return metaFiles, nil
}
