package service

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"htpatcher/internal/domain"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ExportService handles exporting patched files
type ExportService struct {
	logger Logger
}

// NewExportService creates a new export service
func NewExportService(logger Logger) *ExportService {
	return &ExportService{logger: logger}
}

// ExportPatchedFiles exports all patched files to a ZIP archive
func (s *ExportService) ExportPatchedFiles(ctx context.Context, gameDir string, friendlyName string) error {
	// Read patch summary
	summaryPath := filepath.Join(gameDir, "patch-summary.json")
	summaryData, err := os.ReadFile(summaryPath)
	if err != nil {
		s.logger.Error("Failed to read patch-summary.json - game may not be patched")
		return errors.New("patch-summary.json not found - game may not be patched")
	}

	var patchSummary domain.PatchSummary
	if err := json.Unmarshal(summaryData, &patchSummary); err != nil {
		s.logger.Error("Failed to parse patch-summary.json")
		return err
	}

	if len(patchSummary.PatchedFiles) == 0 {
		s.logger.Error("No patched files found in patch summary")
		return errors.New("no patched files found")
	}

	// Sanitize friendly name for filename
	safeName := sanitizeFilename(friendlyName)
	defaultFilename := safeName + "_patched_files.zip"

	// Open save dialog
	outputPath, err := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
		Title:           "Export Patched Files",
		DefaultFilename: defaultFilename,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "ZIP Archive",
				Pattern:     "*.zip",
			},
		},
	})
	if err != nil {
		return err
	}

	// User cancelled
	if outputPath == "" {
		return nil
	}

	// Create ZIP file
	s.logger.Info("Creating ZIP archive...")
	zipFile, err := os.Create(outputPath)
	if err != nil {
		s.logger.Error("Failed to create ZIP file")
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add each patched file to the ZIP
	filesAdded := 0
	for _, relPath := range patchSummary.PatchedFiles {
		srcPath := filepath.Join(gameDir, relPath)

		// Check if file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			s.logger.Info("Skipping missing file: " + relPath)
			continue
		}

		// Open source file
		srcFile, err := os.Open(srcPath)
		if err != nil {
			s.logger.Error("Failed to open file: " + relPath)
			return err
		}

		// Create entry in ZIP (use forward slashes for ZIP compatibility)
		zipPath := filepath.ToSlash(relPath)
		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			srcFile.Close()
			s.logger.Error("Failed to create ZIP entry: " + relPath)
			return err
		}

		// Copy file contents
		_, err = io.Copy(writer, srcFile)
		srcFile.Close()
		if err != nil {
			s.logger.Error("Failed to write file to ZIP: " + relPath)
			return err
		}

		filesAdded++
	}

	s.logger.Success(fmt.Sprintf("Exported %d patched files to ZIP", filesAdded))
	return nil
}

// sanitizeFilename removes or replaces characters that are invalid in filenames
func sanitizeFilename(name string) string {
	// Characters not allowed in Windows filenames
	invalidChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	result := name
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}
	// Trim spaces and dots from ends
	result = strings.Trim(result, " .")
	// If empty, use default
	if result == "" {
		result = "game"
	}
	return result
}
