package service

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/domain/rpgmaker"
	"htpatcher/internal/patcher"
	"htpatcher/internal/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// PatchRepositoryInterface interface with all patch operations
type PatchRepositoryInterface interface {
	Open(path string) (*zip.ReadCloser, error)
	ReadDictionary(zipReader *zip.ReadCloser) (map[string]string, error)
	ReadConfig(zipReader *zip.ReadCloser) (*domain.Config, error)
	GetAllOverrides(zipReader *zip.ReadCloser) ([]string, error)
	ReadFileFromZip(zipReader *zip.ReadCloser, path string) ([]byte, error)
}

// PatchService handles patch operations
type PatchService struct {
	patchRepo      PatchRepositoryInterface
	patcherEngine  *patcher.Engine
	pluginPatcher  *patcher.PluginPatcher
	creditsPatcher *patcher.CreditsPatcher
	logger         Logger
}

// NewPatchService creates a new patch service
func NewPatchService(patchRepo PatchRepositoryInterface, logger Logger) *PatchService {
	return &PatchService{
		patchRepo:      patchRepo,
		patcherEngine:  patcher.NewEngine(logger),
		pluginPatcher:  patcher.NewPluginPatcher(logger),
		creditsPatcher: patcher.NewCreditsPatcher(),
		logger:         logger,
	}
}

// SelectPatchFile opens a file dialog to select a patch file
func (s *PatchService) SelectPatchFile(ctx context.Context) (*domain.PatchInfo, error) {
	filePath, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select the Patch file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Patch file",
				Pattern:     "*.htpatch",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return s.LoadPatchInfo(filePath)
}

// LoadPatchInfo loads patch information from a file
func (s *PatchService) LoadPatchInfo(filePath string) (*domain.PatchInfo, error) {
	patchInfo := &domain.PatchInfo{
		PatchPath: filePath,
	}

	// Open the patch file
	r, err := s.patchRepo.Open(patchInfo.PatchPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Read config
	patchInfo.Config, err = s.patchRepo.ReadConfig(r)
	if err != nil {
		return nil, err
	}

	if patchInfo.Config.Version > Version {
		s.logger.Error(fmt.Sprintf("Patch version %d is not supported.", patchInfo.Config.Version))
		s.logger.Error("Please update the patcher to the latest version.")
		return nil, errors.New("patch version is not supported")
	}

	// Read dictionary
	patchInfo.Dictionary, err = s.patchRepo.ReadDictionary(r)
	if err != nil {
		return nil, err
	}

	// Get overrides
	patchInfo.Overrides, err = s.patchRepo.GetAllOverrides(r)
	if err != nil {
		return nil, err
	}

	return patchInfo, nil
}

// FetchAllPatches fetches all available patches from the API
func (s *PatchService) FetchAllPatches() ([]domain.PatchEntry, error) {
	response, err := http.Get("https://htranslations.com/api/patches")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var patches []domain.PatchEntry
	if err := json.Unmarshal(body, &patches); err != nil {
		return nil, err
	}
	return patches, nil
}

// ApplyPatch applies a patch to a game
func (s *PatchService) ApplyPatch(ctx context.Context, gameInfo *domain.GameInfo, patchInfo *domain.PatchInfo) error {
	s.logger.Info("Starting patch application...")

	// Track all patched files (relative paths from game directory)
	var patchedFiles []string

	// List json files in data folder
	s.logger.Info("Scanning data folder for JSON files...")
	jsonFiles, err := util.ListFilesWithExtension(gameInfo.DataPath, ".json")
	if err != nil {
		s.logger.Error("Failed to scan data folder")
		return err
	}
	s.logger.Info(fmt.Sprintf("Found %d JSON files to patch", len(jsonFiles)))

	// Patch all data files
	for _, jsonFile := range jsonFiles {
		err = s.patcherEngine.PatchDataFile(ctx, jsonFile, patchInfo)
		if err != nil {
			s.logger.Error("Error patching file: " + filepath.Base(jsonFile))
			return err
		}
		// Track patched file
		relPath, _ := filepath.Rel(gameInfo.GameDir, jsonFile)
		patchedFiles = append(patchedFiles, relPath)
	}

	// Patch plugins.js
	pluginsJsPath := filepath.Join(gameInfo.JsPath, "plugins.js")
	err = s.pluginPatcher.UpdatePluginsJs(ctx, pluginsJsPath, patchInfo.Config.PluginsToPatch, patchInfo.Dictionary)
	if err != nil {
		s.logger.Error("Failed to update plugins.js")
		return err
	}

	// Track plugins.js if we have plugins to patch
	if len(patchInfo.Config.PluginsToPatch) > 0 {
		relPath, _ := filepath.Rel(gameInfo.GameDir, pluginsJsPath)
		patchedFiles = append(patchedFiles, relPath)
	}

	// Apply replace rules
	for _, pluginToPatch := range patchInfo.Config.PluginsToPatch {
		for _, replaceRule := range pluginToPatch.ReplaceRules {
			err = s.pluginPatcher.ApplyReplaceRule(ctx, gameInfo.JsPath, pluginToPatch.Plugin, replaceRule)
			if err != nil {
				s.logger.Error("Failed to apply plugin replace rule")
				return err
			}
			// Track patched plugin file
			pluginJsPath := filepath.Join(gameInfo.JsPath, "plugins", pluginToPatch.Plugin+".js")
			relPath, _ := filepath.Rel(gameInfo.GameDir, pluginJsPath)
			// Only add if not already in the list
			found := false
			for _, f := range patchedFiles {
				if f == relPath {
					found = true
					break
				}
			}
			if !found {
				patchedFiles = append(patchedFiles, relPath)
			}
		}
	}

	// Apply overrides
	if len(patchInfo.Overrides) > 0 {
		r, err := s.patchRepo.Open(patchInfo.PatchPath)
		if err != nil {
			s.logger.Error("Failed to open patch")
			return err
		}
		defer r.Close()
		for _, override := range patchInfo.Overrides {
			data, err := s.patchRepo.ReadFileFromZip(r, filepath.Join("overrides", override))
			if err != nil {
				s.logger.Error("Failed to read override")
				return err
			}
			err = os.WriteFile(filepath.Join(gameInfo.GameDir, override), data, 0644)
			if err != nil {
				s.logger.Error("Failed to write override")
				return err
			}
			s.logger.Info(fmt.Sprintf("Overwritten file %s", override))
			// Track override file
			patchedFiles = append(patchedFiles, override)
		}
	}

	// Read system information for credits
	s.logger.Info("Reading system information...")
	systemInfoData, err := os.ReadFile(filepath.Join(gameInfo.DataPath, "system.json"))
	if err != nil {
		s.logger.Error("Failed to read system.json")
		return err
	}

	var systemInfo rpgmaker.System
	if err := json.Unmarshal(systemInfoData, &systemInfo); err != nil {
		s.logger.Error("Failed to parse system.json")
		return err
	}

	// Find main screen image
	s.logger.Info("Looking for main screen image...")
	mainScreenImageName := systemInfo.Title1Name
	pngPath := filepath.Join(gameInfo.ImgPath, "titles1", mainScreenImageName+".png")
	if _, err := os.Stat(pngPath); os.IsNotExist(err) {
		pngPath = filepath.Join(gameInfo.ImgPath, "titles1", mainScreenImageName+".rpgmvp")
	}
	if _, err := os.Stat(pngPath); os.IsNotExist(err) {
		pngPath = filepath.Join(gameInfo.ImgPath, "titles1", mainScreenImageName+".png_")
	}
	if _, err := os.Stat(pngPath); os.IsNotExist(err) {
		s.logger.Error("Main screen image not found")
		return errors.New("main screen image not found")
	}

	// Set default credits location
	if patchInfo.Config.CreditsLocation == "" {
		patchInfo.Config.CreditsLocation = "bottom_left"
	}

	// Add credits to main screen
	s.logger.Info("Adding credits to main screen image...")
	err = s.creditsPatcher.AddCreditsToResource(pngPath, systemInfo.EncryptionKey, patchInfo.Config.CreditsLocation)
	if err != nil {
		s.logger.Error("Failed to add credits")
		return err
	}

	// Track the title image
	titleImageRelPath, _ := filepath.Rel(gameInfo.GameDir, pngPath)
	patchedFiles = append(patchedFiles, titleImageRelPath)

	// Save patch summary
	s.logger.Info("Saving patch summary...")
	patchSummary := domain.PatchSummary{
		PatchedAt:    time.Now().UTC().Format(time.RFC3339),
		PatchedFiles: patchedFiles,
	}
	summaryData, err := json.MarshalIndent(patchSummary, "", "  ")
	if err != nil {
		s.logger.Error("Failed to marshal patch summary")
		return err
	}
	summaryPath := filepath.Join(gameInfo.GameDir, "patch-summary.json")
	if err := os.WriteFile(summaryPath, summaryData, 0644); err != nil {
		s.logger.Error("Failed to write patch summary")
		return err
	}

	s.logger.Success("✓ Patch applied successfully!")
	return nil
}
