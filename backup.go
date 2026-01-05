package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func BackupGameData(a *App, gameInfo GameInfo, patchInfo PatchInfo) error {
	filesToBackup := []string{}

	// List all json files in the data folder
	jsonFiles, err := filepath.Glob(filepath.Join(gameInfo.DataPath, "*.json"))
	if err != nil {
		return err
	}

	// parse system.json
	systemJsonPath := filepath.Join(gameInfo.DataPath, "system.json")
	systemJson, err := os.ReadFile(systemJsonPath)
	if err != nil {
		return err
	}
	var systemInfo System
	json.Unmarshal(systemJson, &systemInfo)

	if systemInfo.Title1Name != "" {
		// glob the file extension .png, .rpgmvp, .png_
		titlesFiles, err := filepath.Glob(filepath.Join(gameInfo.ImgPath, "titles1", systemInfo.Title1Name+".*"))
		if err != nil {
			return err
		}
		for _, titleFile := range titlesFiles {
			relPath, err := filepath.Rel(gameInfo.GameDir, titleFile)
			if err != nil {
				return err
			}
			filesToBackup = append(filesToBackup, relPath)
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

	a.Log(fmt.Sprintf("Found %d files to backup", len(filesToBackup)))

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

	a.Log(fmt.Sprintf("Backed up %d unseen files", copiedFiles))

	return nil
}

func RestoreBackup(a *App, gameInfo GameInfo) error {
	backupPath := filepath.Join(gameInfo.GameDir, ".backup")
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return errors.New("backup folder not found")
	}

	a.Log(fmt.Sprintf("Restoring backup from %s to %s", backupPath, gameInfo.GameDir))

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

		a.Log(fmt.Sprintf("Restored file %s", relPath))

		return nil
	})

	if err != nil {
		return err
	}

	a.Log(fmt.Sprintf("Restored %d files from backup", restoredFiles))

	// Remove the backup folder
	os.RemoveAll(backupPath)
	a.Log(fmt.Sprintf("Removed backup folder %s", backupPath))

	return nil
}

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
