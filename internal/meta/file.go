package meta

import (
	"encoding/json"
	"os"

	"golang.org/x/exp/slog"
)

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
