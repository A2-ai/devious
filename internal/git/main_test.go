package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetRelativePath(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp(".", "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create file
	filePath := filepath.Join(subDir, "test.txt")
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Get relative path
	relativePath, err := GetRelativePath(tempDir, subDir)
	if err != nil {
		t.Error(err)
	}
	if relativePath != "/subdir" && relativePath != "\\subdir" {
		t.Error("Expected relative path to be /subdir or \\subdir")
	}
}

func TestIsGitRepository(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp(".", "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create file
	filePath := filepath.Join(subDir, "test.txt")
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Check if tempDir is a git repository
	isGitRepo := isGitRepository(tempDir)
	if isGitRepo {
		t.Error("Expected tempDir to not be a git repository")
	}

	// Check if subDir is a git repository
	isGitRepo = isGitRepository(subDir)
	if isGitRepo {
		t.Error("Expected subDir to not be a git repository")
	}

	// Create .git
	gitDir := filepath.Join(tempDir, ".git")
	err = os.Mkdir(gitDir, 0755)

	// Check if tempDir is a git repository
	isGitRepo = isGitRepository(tempDir)
	if !isGitRepo {
		t.Error("Expected tempDir to be a git repository")
	}

	// Check if subDir is a git repository
	isGitRepo = isGitRepository(subDir)
	if isGitRepo {
		t.Error("Expected subDir to not be a git repository")
	}
}

func TestGetRootDir(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp(".", "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()
	tempDirAbs, err := filepath.Abs(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create .git
	gitDir := filepath.Join(tempDir, ".git")
	err = os.Mkdir(gitDir, 0755)

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create file
	filePath := filepath.Join(subDir, "test.txt")
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Get repository root
	repoRoot, err := GetNearestRepoDir(subDir)
	if err != nil {
		t.Error(err)
	}
	if repoRoot != tempDirAbs {
		t.Error("Expected repository root to be tempDir, got", repoRoot)
	}

	// Get repository root
	repoRoot2, err := GetNearestRepoDir(filePath)
	if err != nil {
		t.Error(err)
	}
	if repoRoot2 != tempDirAbs {
		t.Error("Expected repository root to be tempDir, got", repoRoot2)
	}
}
