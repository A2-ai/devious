package git

import (
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

func AddToGitIgnore(gitDir string, filePath string) error {
	ignoreFile, err := os.OpenFile(filepath.Join(gitDir, ".gitignore"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer ignoreFile.Close()

	_, err = ignoreFile.WriteString("\n\n# Devious entry\n/" + filepath.Base(filePath))
	if err != nil {
		return err
	}

	slog.Debug("Added file to gitignore", slog.String("git-dir", gitDir), slog.String("file", filePath))

	return err
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
