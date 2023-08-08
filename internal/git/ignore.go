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
	ignoreEntry, err := GetRelativePath(gitDir, path)
	if err != nil {
		return err
	}

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
	if strings.Contains(string(fileBytes), ignoreEntry) {
		slog.Debug("Path already in gitignore file, not adding", slog.String("git-dir", gitDir), slog.String("file", path))
		return nil
	}

	// Add the file to the gitignore file
	_, err = ignoreFile.WriteString("\n\n# Devious entry\n" + ignoreEntry)
	if err != nil {
		slog.Error("Failed to write to gitignore file", slog.String("git-dir", gitDir))
		return err
	}

	slog.Debug("Adding file to gitignore", slog.String("git-dir", gitDir), slog.String("file", path))

	return err
}

func RemoveIgnoreEntry(gitDir string, path string, dry bool) error {
	// Get relative path
	ignoreEntry, err := GetRelativePath(gitDir, path)
	if err != nil {
		return err
	}

	// Open the gitignore file
	// if the gitignore file doesn't exist, there's nothing to do
	ignoreFilePath := filepath.Join(gitDir, ".gitignore")
	if _, err := os.Stat(ignoreFilePath); os.IsNotExist(err) {
		slog.Debug("gitignore file does not exist, nothing to do", slog.String("git-dir", gitDir))
		return nil
	}
	ignoreFile, err := os.OpenFile(ignoreFilePath, os.O_RDWR, 0644)
	if err != nil {
		slog.Error("Failed to open gitignore file", slog.String("git-dir", gitDir))
		return err
	}
	defer ignoreFile.Close()

	// Find starting line to remove
	ignoreBytes, err := io.ReadAll(ignoreFile)
	ignoreFile.Seek(0, 0)
	if err != nil {
		slog.Error("Failed to read gitignore file", slog.String("git-dir", gitDir))
		return err
	}
	ignoreContents := strings.Split(string(ignoreBytes), "\n")
	for i, line := range ignoreContents {
		if line == ignoreEntry {
			err := removeLines(ignoreFile, i-1, 3)
			if err != nil {
				slog.Error("Failed to remove lines from gitignore file", slog.String("error", err.Error()))
				return err
			}
		}
	}

	if dry {
		slog.Debug("Dry run: removed gitignore entry", slog.String("git-dir", gitDir), slog.String("file", path))
		return nil
	}

	slog.Debug("Removed gitignore entry", slog.String("git-dir", gitDir), slog.String("file", path))

	return nil
}
