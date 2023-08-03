package storage

import (
	"devious/internal/config"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

// Copies a file from the source path to the destination path
func Copy(srcPath string, destPath string, conf config.Config) error {
	// Open source file
	src, err := os.Open(srcPath)
	if err == os.ErrNotExist {
		slog.Error("File does not exist", slog.String("path", srcPath))
		return err
	} else if err != nil {
		slog.Error("Failed to open file", slog.String("path", srcPath))
		return err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(destPath)

	// Create the directory if it doesn't exist
	// Return if there was an error other than the directory not existing
	if err == os.ErrNotExist {
		err = os.MkdirAll(filepath.Dir(destPath), 0755)
		if err != nil {
			slog.Error("Failed to create directory", slog.String("path", filepath.Dir(destPath)))
			return err
		}
	} else if err != nil {
		slog.Error("Failed to create file", slog.String("path", destPath))
		return err
	}

	defer dst.Close()

	// Copy the file
	slog.Info("Copying file...")
	copiedBytes, err := io.Copy(dst, src)
	if err != nil {
		slog.Error("Failed to copy file", slog.String("path", srcPath))
		return err
	}
	slog.Info("Copied file",
		slog.String("filesize", fmt.Sprintf("%1.2f MB", float64(copiedBytes)/1000000)),
		slog.String("from", srcPath), slog.String("to", destPath))

	return nil
}
