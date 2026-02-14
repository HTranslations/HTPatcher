package rgss3a

import (
	"fmt"
	"os"
	"path/filepath"
)

// ExtractToDir extracts all files from the archive to the specified directory.
func ExtractToDir(archivePath, outputDir string) error {
	archive, err := Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer archive.Close()

	files := archive.List()
	for _, filePath := range files {
		data, err := archive.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", filePath, err)
		}

		// Construct output path
		outPath := filepath.Join(outputDir, filePath)

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
		}

		// Write file
		if err := os.WriteFile(outPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filePath, err)
		}
	}

	return nil
}

// ExtractFile extracts a single file from the archive to the specified path.
func ExtractFile(archivePath, filePath, outputPath string) error {
	archive, err := Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer archive.Close()

	data, err := archive.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
