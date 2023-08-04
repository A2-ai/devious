package git

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slog"
)

func AddToIgnore(gitDir string, path string, dry bool) error {
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
