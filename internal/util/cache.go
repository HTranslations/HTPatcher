package util

import (
	"os"
	"path/filepath"
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

	return filepath.Join(cacheDir, "htpatcher_update.exe"), nil
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
