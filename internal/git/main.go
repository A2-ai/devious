package git

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slog"
)

// Gets a file path relative to the root
func GetRelativePath(rootDir string, filePath string) (string, error) {
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		slog.Error("Failed to get absolute path for root dir", slog.String("path", filePath))
		return "", err
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		slog.Error("Failed to get absolute path", slog.String("path", filePath))
		return "", err
	}

	return strings.TrimPrefix(absFilePath, absRootDir), nil
}

// Checks if the supplied path
func isGitRepository(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// Recursively gets the nearest git repository root, given an absolute path
func getNearestRepoDir(dir string) (string, error) {
	// Check if the current directory is a git repository
	if isGitRepository(dir) {
		return dir, nil
	}

	// Get parent directory
	parentDir := filepath.Dir(dir)

	// Stop if we're at the root directory
	if parentDir == "/" {
		return "", os.ErrNotExist
	}

	// Recursively call this function with the parent directory
	return getNearestRepoDir(parentDir)
}

// Gets the nearest parent repository root of the supplied directory, or an error if the directory is not contained within a git repository
func GetNearestRepoDir(dir string) (string, error) {
	// Get the absolute path of the supplied directory
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	// Get the nearest repository root
	repoRoot, err := getNearestRepoDir(absDir)
	return repoRoot, err
}
