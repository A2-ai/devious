package git

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slog"
)

// Gets a file path relative to the repository root
func GetRelativePath(repoRoot string, filePath string) string {
	return strings.TrimPrefix(filePath, repoRoot)
}

// Checks if the supplied path
func isGitRepository(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// Gets the nearest repository root of the supplied dir, or an error if the dir is not contained within a git repository
func getRepositoryRoot(dir string) (string, error) {
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
	return getRepositoryRoot(parentDir)
}

func GetRootDir() (string, error) {
	// Start at the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Get the repository root
	repoRoot, err := getRepositoryRoot(cwd)
	if err != nil {
		return "", err
	}

	slog.Debug("Found git repository at", slog.String("root", repoRoot))

	return repoRoot, nil
}
