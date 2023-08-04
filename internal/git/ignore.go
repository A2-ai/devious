package git

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slog"
)

func AddIgnoreEntry(gitDir string, path string, dry bool) error {
	// Get relative path
	relativePath := GetRelativePath(gitDir, path)

	// Open the gitignore file, creating one if it doesn't exist
	ignoreFile, err := os.OpenFile(filepath.Join(gitDir, ".gitignore"), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		slog.Error("Failed to open gitignore file", slog.String("git-dir", gitDir))
		return err
	}
	defer ignoreFile.Close()

	if dry {
		slog.Debug("Dry run: added file to gitignore", slog.String("git-dir", gitDir), slog.String("file", path))
		return nil
	}

	// Check if the file is already in the gitignore file
	fileBytes, err := io.ReadAll(ignoreFile)
	if err != nil {
		slog.Error("Failed to read gitignore file", slog.String("git-dir", gitDir))
		return err
	}
	if strings.Contains(string(fileBytes), relativePath) {
		slog.Debug("Path already in gitignore file, not adding", slog.String("git-dir", gitDir), slog.String("file", path))
		return nil
	}

	// Add the file to the gitignore file
	_, err = ignoreFile.WriteString("\n\n# Devious entry\n" + relativePath)
	if err != nil {
		slog.Error("Failed to write to gitignore file", slog.String("git-dir", gitDir))
		return err
	}

	slog.Debug("Adding file to gitignore", slog.String("git-dir", gitDir), slog.String("file", path))

	return err
}

func RemoveIgnoreEntry(gitDir string, path string, dry bool) error {
	// Get relative path
	// relativePath := GetRelativePath(gitDir, path)

	// Open the gitignore file
	// if the gitignore file doesn't exist, there's nothing to do
	ignoreFilePath := filepath.Join(gitDir, ".gitignore")
	if _, err := os.Stat(ignoreFilePath); os.IsNotExist(err) {
		return nil
	}
	ignoreFile, err := os.OpenFile(ignoreFilePath, os.O_RDWR, 0644)
	if err != nil {
		slog.Error("Failed to open gitignore file", slog.String("git-dir", gitDir))
		return err
	}
	defer ignoreFile.Close()

	if dry {
		slog.Debug("Dry run: removed path from gitignore", slog.String("git-dir", gitDir), slog.String("file", path))
		return nil
	}

	return nil
}
