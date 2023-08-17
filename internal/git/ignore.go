package git

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func AddIgnoreEntry(gitDir string, path string, dry bool) error {
	// Get relative path
	ignoreEntry, err := GetRelativePath(gitDir, path)
	if err != nil {
		return err
	}
	ignoreEntry = "/" + ignoreEntry // Add leading slash

	// Open the gitignore file, creating one if it doesn't exist
	ignoreFile, err := os.OpenFile(filepath.Join(gitDir, ".gitignore"), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer ignoreFile.Close()

	if dry {
		return nil
	}

	// Check if the file is already in the gitignore file
	fileBytes, err := io.ReadAll(ignoreFile)
	if err != nil {
		return err
	}
	if strings.Contains(string(fileBytes), ignoreEntry) {
		return nil
	}

	// Add the file to the gitignore file
	_, err = ignoreFile.WriteString("\n\n# Devious entry\n" + ignoreEntry)
	if err != nil {
		return err
	}

	return err
}

func RemoveIgnoreEntry(gitDir string, path string, dry bool) error {
	// Get relative path
	ignoreEntry, err := GetRelativePath(gitDir, path)
	if err != nil {
		return err
	}
	ignoreEntry = "/" + ignoreEntry // Add leading slash

	// Open the gitignore file
	// if the gitignore file doesn't exist, there's nothing to do
	ignoreFilePath := filepath.Join(gitDir, ".gitignore")
	if _, err := os.Stat(ignoreFilePath); os.IsNotExist(err) {
		return nil
	}
	ignoreFile, err := os.OpenFile(ignoreFilePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer ignoreFile.Close()

	// Find starting line to remove
	ignoreBytes, err := io.ReadAll(ignoreFile)
	ignoreFile.Seek(0, 0)
	if err != nil {
		return err
	}
	ignoreContents := strings.Split(string(ignoreBytes), "\n")
	for i, line := range ignoreContents {
		if line == ignoreEntry {
			err := removeLines(ignoreFile, i-1, 3)
			if err != nil {
				return err
			}
		}
	}

	if dry {
		return nil
	}

	return nil
}
