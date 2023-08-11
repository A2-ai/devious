package meta

import (
	"encoding/json"
	"os"

	"golang.org/x/exp/slog"
)

// Creates a metadata file
func Save(metadata Metadata, path string, dry bool) (err error) {
	var metadataFile *os.File
	if dry {
		return nil
	}

	metadataFile, err = os.Create(path + FileExtension)
	if err != nil {
		slog.Error("Failed to create metadata file", slog.String("path", path))
		return err
	}
	defer metadataFile.Close()

	jsonBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		slog.Error("Failed to encode metadata", slog.String("path", path))
		return err
	}

	_, err = metadataFile.Write(jsonBytes)
	if err != nil {
		slog.Error("Failed to write metadata", slog.String("path", path))
		return err
	}

	return nil
}

// Loads a metadata file
func Load(path string) (metadata Metadata, err error) {
	metadataFile, err := os.Open(path + FileExtension)
	if err != nil {
		return metadata, err
	}

	err = json.NewDecoder(metadataFile).Decode(&metadata)
	if err != nil {
		return metadata, err
	}

	return metadata, nil
}
