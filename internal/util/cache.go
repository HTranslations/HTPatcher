package util

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetUpdateCacheDir returns the update cache directory path
// Windows: C:\Users\<user>\AppData\Local\htpatcher\updates
func GetUpdateCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	updateCacheDir := filepath.Join(cacheDir, "htpatcher", "updates")

	// Ensure directory exists
	if err := os.MkdirAll(updateCacheDir, 0755); err != nil {
		return "", err
	}

	return updateCacheDir, nil
}

// GetUpdateExePath returns the full path to the downloaded update executable
func GetUpdateExePath() (string, error) {
	cacheDir, err := GetUpdateCacheDir()
	if err != nil {
		return "", err
	}

	filename := "htpatcher_update"
	if runtime.GOOS == "windows" {
		filename = "htpatcher_update.exe"
	}

	return filepath.Join(cacheDir, filename), nil
}

// CleanUpdateCache deletes all files in the update cache directory
func CleanUpdateCache() error {
	cacheDir, err := GetUpdateCacheDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if err := os.Remove(filepath.Join(cacheDir, entry.Name())); err != nil {
			// Continue even if some files fail to delete
			continue
		}
	}

	return nil
}
