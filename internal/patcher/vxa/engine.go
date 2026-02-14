package vxa

import (
	"context"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/marshal"
	"htpatcher/internal/rgss3a"
	"htpatcher/internal/util"
	"os"
	"path/filepath"
	"strings"
)

// Logger interface for logging operations
type Logger interface {
	Info(message string)
	Success(message string)
	Error(message string)
	Warn(message string)
}

// Engine handles VX Ace patching operations
type Engine struct {
	logger Logger
}

// NewEngine creates a new VX Ace patcher engine
func NewEngine(logger Logger) *Engine {
	return &Engine{
		logger: logger,
	}
}

// PatchGame patches a VX Ace game with translations
func (e *Engine) PatchGame(ctx context.Context, gameInfo *domain.GameInfo, patchInfo *domain.PatchInfo) ([]string, error) {
	var patchedFiles []string

	// Check if we need to extract from archive first
	rgss3aPath := filepath.Join(gameInfo.GameDir, "Game.rgss3a")
	if _, err := os.Stat(rgss3aPath); err == nil {
		e.logger.Info("Extracting Game.rgss3a archive...")
		if err := e.extractArchive(rgss3aPath, gameInfo.GameDir); err != nil {
			return nil, fmt.Errorf("failed to extract archive: %w", err)
		}
		e.logger.Success("Archive extracted successfully")

		// Delete the archive so game reads from extracted files
		if err := os.Remove(rgss3aPath); err != nil {
			e.logger.Warn(fmt.Sprintf("Failed to remove archive: %v", err))
		} else {
			e.logger.Info("Removed Game.rgss3a (game will now read from Data folder)")
		}
	}

	// Scan for .rvdata2 files
	dataPath := gameInfo.DataPath
	if dataPath == "" {
		dataPath = filepath.Join(gameInfo.GameDir, "Data")
	}

	e.logger.Info(fmt.Sprintf("Scanning %s for .rvdata2 files...", dataPath))

	rvdata2Files, err := util.ListFilesWithExtension(dataPath, ".rvdata2")
	if err != nil {
		return nil, fmt.Errorf("failed to scan data folder: %w", err)
	}

	e.logger.Info(fmt.Sprintf("Found %d .rvdata2 files to patch", len(rvdata2Files)))

	// Patch each file
	for _, filePath := range rvdata2Files {
		if err := e.patchDataFile(ctx, filePath, patchInfo); err != nil {
			e.logger.Error(fmt.Sprintf("Error patching %s: %v", filepath.Base(filePath), err))
			return nil, err
		}
		relPath, _ := filepath.Rel(gameInfo.GameDir, filePath)
		patchedFiles = append(patchedFiles, relPath)
	}

	return patchedFiles, nil
}

// extractArchive extracts a Game.rgss3a archive
func (e *Engine) extractArchive(archivePath, outputDir string) error {
	return rgss3a.ExtractToDir(archivePath, outputDir)
}

// patchDataFile patches a single .rvdata2 file based on its type
func (e *Engine) patchDataFile(ctx context.Context, filePath string, patchInfo *domain.PatchInfo) error {
	filename := filepath.Base(filePath)
	e.logger.Info("Patching: " + filename)

	fileType := e.getDataFileType(filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var patchedData []byte

	switch fileType {
	case "actors":
		patchedData, err = patchActors(data, patchInfo)
	case "classes":
		patchedData, err = patchClasses(data, patchInfo)
	case "skills":
		patchedData, err = patchSkills(data, patchInfo)
	case "items":
		patchedData, err = patchItems(data, patchInfo)
	case "weapons":
		patchedData, err = patchWeapons(data, patchInfo)
	case "armors":
		patchedData, err = patchArmors(data, patchInfo)
	case "enemies":
		patchedData, err = patchEnemies(data, patchInfo)
	case "states":
		patchedData, err = patchStates(data, patchInfo)
	case "troops":
		patchedData, err = patchTroops(data, patchInfo)
	case "commonevents":
		patchedData, err = patchCommonEvents(data, patchInfo)
	case "map":
		patchedData, err = patchMap(data, patchInfo)
	case "system":
		patchedData, err = patchSystem(data, patchInfo)
	case "scripts":
		patchedData, err = patchScripts(data, patchInfo)
	default:
		// Unknown file type, skip
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to patch %s: %w", filename, err)
	}

	if patchedData == nil {
		return nil
	}

	if err := os.WriteFile(filePath, patchedData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetSystemTitleImageName reads System.rvdata2 and returns the title1_name property
func (e *Engine) GetSystemTitleImageName(dataPath string) (string, error) {
	data, err := os.ReadFile(filepath.Join(dataPath, "System.rvdata2"))
	if err != nil {
		return "", fmt.Errorf("failed to read System.rvdata2: %w", err)
	}

	raw, err := marshal.Parse(data)
	if err != nil {
		return "", fmt.Errorf("failed to parse System.rvdata2: %w", err)
	}

	obj, ok := raw.(*marshal.RubyObject)
	if !ok {
		return "", fmt.Errorf("expected RubyObject, got %T", raw)
	}

	title1Name, ok := obj.Properties["title1_name"].(string)
	if !ok {
		return "", fmt.Errorf("title1_name not found in System.rvdata2")
	}

	return title1Name, nil
}

// getDataFileType determines the file type from the filename
func (e *Engine) getDataFileType(filePath string) string {
	filename := filepath.Base(filePath)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	filename = strings.ToLower(filename)

	fileTypeMap := map[string]string{
		"actors":       "actors",
		"classes":      "classes",
		"skills":       "skills",
		"items":        "items",
		"weapons":      "weapons",
		"armors":       "armors",
		"enemies":      "enemies",
		"states":       "states",
		"troops":       "troops",
		"commonevents": "commonevents",
		"system":       "system",
		"scripts":      "scripts",
		"mapinfos":     "mapinfos",
		"animations":   "animations",
		"tilesets":     "tilesets",
	}

	if fileType, ok := fileTypeMap[filename]; ok {
		return fileType
	}

	// Check for Map files (Map001.rvdata2, etc.)
	if strings.HasPrefix(filename, "map") && len(filename) > 3 {
		return "map"
	}

	return filename
}
