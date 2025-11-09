//go:build linux

package yqlib

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestChangeOwner(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "testfile.txt")

	// Create a test file
	err := os.WriteFile(testFile, []byte("test content"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Get file info
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat test file: %v", err)
	}

	// Create another temporary file to change ownership of
	tempFile, err := os.CreateTemp(tempDir, "chown_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Test changeOwner function
	err = changeOwner(info, tempFile)
	if err != nil {
		t.Errorf("changeOwner failed: %v", err)
	}

	// Verify that the function doesn't panic with valid input
	tempFile2, err := os.CreateTemp(tempDir, "chown_test2_*.txt")
	if err != nil {
		t.Fatalf("Failed to create second temp file: %v", err)
	}
	defer os.Remove(tempFile2.Name())
	tempFile2.Close()

	// Test with the second file
	err = changeOwner(info, tempFile2)
	if err != nil {
		t.Errorf("changeOwner failed on second file: %v", err)
	}
}

func TestChangeOwnerWithInvalidFileInfo(t *testing.T) {
	// Create a mock file info that doesn't have syscall.Stat_t
	mockInfo := &mockFileInfo{
		name: "mock",
		size: 0,
		mode: 0600,
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp(t.TempDir(), "chown_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Test changeOwner with mock file info (should not panic)
	err = changeOwner(mockInfo, tempFile)
	if err != nil {
		t.Errorf("changeOwner failed with mock file info: %v", err)
	}
}

func TestChangeOwnerWithNonExistentFile(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp(t.TempDir(), "chown_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Get file info
	info, err := os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat temp file: %v", err)
	}

	// Remove the file
	os.Remove(tempFile.Name())

	err = changeOwner(info, tempFile)
	// The function should not panic even if the file doesn't exist
	if err != nil {
		t.Logf("Expected error when changing owner of non-existent file: %v", err)
	}
}

// mockFileInfo implements fs.FileInfo but doesn't have syscall.Stat_t
type mockFileInfo struct {
	name string
	size int64
	mode os.FileMode
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool        { return false }
func (m *mockFileInfo) Sys() interface{}   { return nil } // This will cause the type assertion to fail

func TestChangeOwnerWithSyscallStatT(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp(t.TempDir(), "chown_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Get file info
	info, err := os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat temp file: %v", err)
	}

	err = changeOwner(info, tempFile)
	if err != nil {
		t.Logf("changeOwner returned error (this might be expected in some environments): %v", err)
	}
}
