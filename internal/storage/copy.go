package storage

import (
	"dvs/internal/config"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

// Copies a file from the source path to the destination path
func Copy(srcPath string, destPath string, conf config.Config, dry bool) error {
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
	var dst *os.File
	if !dry {
		dst, err = os.Create(destPath)
	}

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

	// Calculate file size in MB
	srcStat, err := src.Stat()
	if err != nil {
		slog.Error("Failed to get file info", slog.String("path", srcPath))
		return err
	}
	fileSize := float64(srcStat.Size()) / 1000000 // 1 MB = 1000000 bytes

	// Copy the file
	if !dry {
		slog.Info("Copying file...")
		_, err := io.Copy(dst, src)
		if err != nil {
			slog.Error("Failed to copy file", slog.String("path", srcPath))
			return err
		}
	} else {
		slog.Info("Dry run: copying file...")
	}

	slog.Info("Copied file",
		slog.String("from", srcPath),
		slog.String("to", destPath),
		slog.String("filesize", fmt.Sprintf("%1.2f MB", fileSize)))

	return nil
}
