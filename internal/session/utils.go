package session

import (
	"os"
	"path/filepath"
	"time"
)

// findNewestFile returns the newest file matching the pattern by modification time
func findNewestFile(pattern string) string {
	files, _ := filepath.Glob(pattern)
	if len(files) == 0 {
		return ""
	}
	if len(files) == 1 {
		return files[0]
	}

	var newestFile string
	var newestTime time.Time

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		if info.ModTime().After(newestTime) {
			newestTime = info.ModTime()
			newestFile = file
		}
	}
	return newestFile
}
