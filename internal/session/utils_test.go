package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFindNewestFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a few files with different modification times
	file1 := filepath.Join(tmpDir, "file1.json")
	file2 := filepath.Join(tmpDir, "file2.json")
	file3 := filepath.Join(tmpDir, "other.txt")

	if err := os.WriteFile(file1, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	// Wait a bit to ensure distinct modification times
	time.Sleep(10 * time.Millisecond)

	if err := os.WriteFile(file2, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(file3, []byte("not matching pattern"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test finding the newest .json file
	pattern := filepath.Join(tmpDir, "*.json")
	newest := findNewestFile(pattern)

	if newest != file2 {
		t.Errorf("Expected newest file to be %s, got %s", file2, newest)
	}

	// Test with no matches
	newest = findNewestFile(filepath.Join(tmpDir, "*.missing"))
	if newest != "" {
		t.Errorf("Expected empty string for no matches, got %s", newest)
	}
}
