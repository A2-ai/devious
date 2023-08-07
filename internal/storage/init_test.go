package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitNoPerms(t *testing.T) {
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

	// Create storage directory with limited permissions
	err = os.Mkdir(filepath.Join(tempDir, "storage"), 0000)
	if err != nil {
		t.Fatal(err)
	}

	// Run function
	err = Init(
		filepath.Join(tempDir, "storage"),
		tempDir,
	)
	if err == nil {
		t.Error("Function did not return error")
	}
}

func TestInitNoParent(t *testing.T) {
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

	// Run function
	err = Init(
		filepath.Join(tempDir, "nested1/nested2"),
		tempDir,
	)
	if err != nil {
		t.Error("Function returned error", err)
	}
}

func TestInitAlreadyExists(t *testing.T) {
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

	// Run function
	err = Init(
		tempDir,
		filepath.Join(tempDir, "storage"),
	)

	// Run function again
	err = Init(
		tempDir,
		filepath.Join(tempDir, "storage"),
	)
	if err != nil {
		t.Error("Function returned error", err)
	}
}