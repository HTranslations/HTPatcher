package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/domain/rpgmaker"
	"htpatcher/internal/util"
	"io"
	"os"
	"path/filepath"
	"slices"
)

// BackupService handles backup and restore operations
type BackupService struct {
	logger Logger
}

// NewBackupService creates a new backup service
func NewBackupService(logger Logger) *BackupService {
	return &BackupService{logger: logger}
}

// BackupGameData creates a backup of game data before patching
func (s *BackupService) BackupGameData(gameInfo *domain.GameInfo, patchInfo *domain.PatchInfo) error {
	filesToBackup := []string{}

	// List all json files in the data folder
	jsonFiles, err := util.ListFilesWithExtension(gameInfo.DataPath, ".json")
	if err != nil {
		return err
	}

	// Parse system.json to get title image
	systemJsonPath := filepath.Join(gameInfo.DataPath, "system.json")
	systemJson, err := os.ReadFile(systemJsonPath)
	if err != nil {
		return err
	}
	var systemInfo rpgmaker.System
	json.Unmarshal(systemJson, &systemInfo)

	if systemInfo.Title1Name != "" {
		// Glob the file extension .png, .rpgmvp, .png_
		titlesPath := filepath.Join(gameInfo.ImgPath, "titles1")
		filesInTitlesPath, err := os.ReadDir(titlesPath)
		if err != nil {
			return err
		}
		possibleFilenames := []string{systemInfo.Title1Name + ".png", systemInfo.Title1Name + ".rpgmvp", systemInfo.Title1Name + ".png_"}
		for _, file := range filesInTitlesPath {
			if slices.Contains(possibleFilenames, file.Name()) {
				relPath, err := filepath.Rel(gameInfo.GameDir, filepath.Join(titlesPath, file.Name()))
				if err != nil {
					return err
				}
				filesToBackup = append(filesToBackup, relPath)
			}
		}
	}

	for _, jsonFile := range jsonFiles {
		relPath, err := filepath.Rel(gameInfo.GameDir, jsonFile)
		if err != nil {
			return err
		}
		filesToBackup = append(filesToBackup, relPath)
	}

	if len(patchInfo.Config.PluginsToPatch) > 0 {
		pluginsJsPath := filepath.Join(gameInfo.JsPath, "plugins.js")
		pluginsJsRelPath, err := filepath.Rel(gameInfo.GameDir, pluginsJsPath)
		if err != nil {
			return err
		}
		filesToBackup = append(filesToBackup, pluginsJsRelPath)
		for _, pluginToPatch := range patchInfo.Config.PluginsToPatch {
			pluginJsPath := filepath.Join(gameInfo.JsPath, "plugins", pluginToPatch.Plugin+".js")
			pluginJsRelPath, err := filepath.Rel(gameInfo.GameDir, pluginJsPath)
			if err != nil {
				return err
			}
			filesToBackup = append(filesToBackup, pluginJsRelPath)
		}
	}

	filesToBackup = append(filesToBackup, patchInfo.Overrides...)

	s.logger.Info(fmt.Sprintf("Found %d files to backup", len(filesToBackup)))

	backupPath := filepath.Join(gameInfo.GameDir, ".backup")
	os.MkdirAll(backupPath, 0755)

	copiedFiles := 0
	for _, file := range filesToBackup {
		srcPath := filepath.Join(gameInfo.GameDir, file)
		dstPath := filepath.Join(backupPath, file)
		copied, err := copyFile(srcPath, dstPath, false)
		if err != nil {
			return err
		}
		if copied {
			copiedFiles++
		}
	}

	s.logger.Info(fmt.Sprintf("Backed up %d unseen files", copiedFiles))
	return nil
}

// RestoreBackup restores a backup
func (s *BackupService) RestoreBackup(gameInfo *domain.GameInfo) error {
	backupPath := filepath.Join(gameInfo.GameDir, ".backup")
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return errors.New("backup folder not found")
	}

	s.logger.Info(fmt.Sprintf("Restoring backup from %s to %s", backupPath, gameInfo.GameDir))

	// Copy recursively from backup path to game path
	restoredFiles := 0
	err := filepath.Walk(backupPath, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the backup directory itself
		if srcPath == backupPath {
			return nil
		}

		// Get relative path from backup folder
		relPath, err := filepath.Rel(backupPath, srcPath)
		if err != nil {
			return err
		}

		// Construct destination path
		dstPath := filepath.Join(gameInfo.GameDir, relPath)

		// If it's a directory, create it
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		// If it's a file, copy it with overwrite
		copied, err := copyFile(srcPath, dstPath, true)
		if err != nil {
			return err
		}
		if copied {
			restoredFiles++
		}

		s.logger.Info(fmt.Sprintf("Restored file %s", relPath))
		return nil
	})

	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("Restored %d files from backup", restoredFiles))

	// Remove the backup folder
	os.RemoveAll(backupPath)
	s.logger.Info(fmt.Sprintf("Removed backup folder %s", backupPath))

	// Remove patch-summary.json since the game is no longer patched
	patchSummaryPath := filepath.Join(gameInfo.GameDir, "patch-summary.json")
	if _, err := os.Stat(patchSummaryPath); err == nil {
		os.Remove(patchSummaryPath)
		s.logger.Info("Removed patch-summary.json")
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(srcPath string, dstPath string, overwrite bool) (bool, error) {
	if _, err := os.Stat(dstPath); err == nil && !overwrite {
		return false, nil
	}
	src, err := os.Open(srcPath)
	if err != nil {
		return false, err
	}
	defer src.Close()
	os.MkdirAll(filepath.Dir(dstPath), 0755)
	dst, err := os.Create(dstPath)
	if err != nil {
		return false, err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return true, err
}
