package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopySanity(t *testing.T) {
	// Make temp directory
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

	// Create file
	filePath := filepath.Join(tempDir, "test1")
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test copying
	err = Copy(filePath, filepath.Join(tempDir, "test2"), false)
	if err != nil {
		t.Error(err)
	}
}

func TestCopyDstAlreadyExists(t *testing.T) {
	// Make temp directory
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

	// Create src file
	srcFilePath := filepath.Join(tempDir, "test1")
	err = os.WriteFile(srcFilePath, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create dst file
	dstFilePath := filepath.Join(tempDir, "test2")
	err = os.WriteFile(dstFilePath, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test copying
	err = Copy(srcFilePath, dstFilePath, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCopyInvalidSrcPath(t *testing.T) {
	// Make temp directory
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

	// Test copying
	err = Copy("invalid.txt", filepath.Join(tempDir, "test.txt"), false)
	if err == nil {
		t.Error("Expected error when copying invalid file")
	}
}

func TestCopyInvalidDstPath(t *testing.T) {
	// Make temp directory
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

	// Create file
	filePath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test copying
	err = Copy(filePath, filepath.Join(tempDir, "invalidDir/test.txt"), false)
	if err != nil {
		t.Error(err)
	}
}

func TestCopyNoSrcPerms(t *testing.T) {
	// Make temp directory
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

	// Create file
	filePath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(filePath, []byte("test"), 000)
	if err != nil {
		t.Fatal(err)
	}

	// Test copying
	err = Copy(filePath, filepath.Join(tempDir, "invalidDir/test.txt"), false)
	if err == nil {
		t.Error("Expected error when copying file with no permissions")
	}
}

func TestCopyNoDstPerms(t *testing.T) {
	// Make temp directory
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

	// Create file
	filePath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(filePath, []byte("test"), 000)
	if err != nil {
		t.Fatal(err)
	}

	// Create dest dir
	err = os.MkdirAll(filepath.Join(tempDir, "dir"), 000)
	if err != nil {
		t.Fatal(err)
	}

	// Test copying
	err = Copy(filePath, filepath.Join(tempDir, "dir/test.txt"), false)
	if err == nil {
		t.Error("Expected error when copying file with no permissions")
	}
}
