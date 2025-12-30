package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type GameInfo struct {
	GameDir   string `json:"gameDir"`
	ExePath   string `json:"exePath"`
	DataPath  string `json:"dataPath"`
	JsPath    string `json:"jsPath"`
	GameTitle string `json:"gameTitle"`
}

type PatchInfo struct {
	PatchPath  string            `json:"patchPath"`
	Dictionary map[string]string `json:"dictionary"`
	Config     *Config           `json:"config"`
}

type PatchEntry struct {
	Title           string `json:"title"`
	SystemGameTitle string `json:"systemGameTitle"`
	PatchDownloadId string `json:"patchDownloadId"`
}

var version = 4

func (a *App) SelectGameExeFile() (*GameInfo, error) {
	gameInfo := GameInfo{}
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select the Game.exe file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Game.exe",
				Pattern:     "*.exe",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Set game paths
	gameInfo.ExePath = filePath
	gameInfo.GameDir = filepath.Dir(filePath)
	a.Log(fmt.Sprintf("Game directory: %s", gameInfo.GameDir))

	// Set data and js paths
	dataPath := filepath.Join(gameInfo.GameDir, "data")
	jsPath := filepath.Join(gameInfo.GameDir, "js")
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		dataPath = filepath.Join(gameInfo.GameDir, "www", "data")
		jsPath = filepath.Join(gameInfo.GameDir, "www", "js")
	}

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		a.LogError("Data directory not found")
		return nil, errors.New("data directory not found")
	}
	if _, err := os.Stat(jsPath); os.IsNotExist(err) {
		a.LogError("JS directory not found")
		return nil, errors.New("js directory not found")
	}

	gameInfo.DataPath = dataPath
	gameInfo.JsPath = jsPath

	systemInfoData, err := os.ReadFile(filepath.Join(gameInfo.DataPath, "system.json"))
	if err != nil {
		a.LogError("Failed to read system.json")
		return nil, err
	}

	var systemInfo System
	if err := json.Unmarshal(systemInfoData, &systemInfo); err != nil {
		a.LogError("Failed to parse system.json")
		return nil, err
	}

	gameInfo.GameTitle = systemInfo.GameTitle
	a.Log(fmt.Sprintf("Game title: \"%s\"", gameInfo.GameTitle))

	return &gameInfo, nil
}

func (a *App) SelectPatchFile() (*PatchInfo, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
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

	patchInfo, err := LoadPatchInfo(a, filePath)
	if err != nil {
		return nil, err
	}
	return patchInfo, nil
}

func LoadPatchInfo(a *App, filePath string) (*PatchInfo, error) {
	patchInfo := &PatchInfo{
		PatchPath: filePath,
	}

	// Set dictionary
	r, err := OpenPatch(patchInfo.PatchPath)
	if err != nil {
		return nil, err
	}
	patchInfo.Dictionary, err = ReadDictionary(r)
	if err != nil {
		return nil, err
	}

	// Set config
	patchInfo.Config, err = ReadConfig(r)
	if err != nil {
		return nil, err
	}

	if patchInfo.Config.Version > version {
		a.LogError(fmt.Sprintf("Patch version %d is not supported.", patchInfo.Config.Version))
		a.LogError("Please update the patcher to the latest version.")
		a.LogError("You can download the latest version from the website.")
		a.LogError("https://htranslations.com")
		return nil, errors.New("patch version is not supported")
	}

	return patchInfo, nil
}

type LogMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"` // "info", "success", "error"
}

func (a *App) Log(message string) {
	runtime.EventsEmit(a.ctx, "log", LogMessage{
		Message: message,
		Type:    "info",
	})
}

func (a *App) LogSuccess(message string) {
	runtime.EventsEmit(a.ctx, "log", LogMessage{
		Message: message,
		Type:    "success",
	})
}

func (a *App) LogError(message string) {
	runtime.EventsEmit(a.ctx, "log", LogMessage{
		Message: message,
		Type:    "error",
	})
}

func (a *App) DownloadPatch(patchDownloadId string) (*PatchInfo, error) {
	a.Log(fmt.Sprintf("Downloading patch with download ID %s", patchDownloadId))
	filePath, err := DownloadPatch(patchDownloadId)
	if err != nil {
		return nil, err
	}

	a.Log(fmt.Sprintf("Downloaded patch to %s", filePath))
	a.Log("Loading patch into memory...")

	patchInfo, err := LoadPatchInfo(a, filePath)
	if err != nil {
		return nil, err
	}

	os.Remove(filePath)
	a.Log(fmt.Sprintf("Removed %s", filePath))

	return patchInfo, nil
}

func (a *App) FetchAllPatches() ([]PatchEntry, error) {
	response, err := http.Get("https://htranslations.com/api/patches")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var patches []PatchEntry
	if err := json.Unmarshal(body, &patches); err != nil {
		return nil, err
	}
	return patches, nil
}

func (a *App) ApplyPatch(gameInfo GameInfo, patchInfo PatchInfo) error {
	a.Log("Starting patch application...")

	// list json files in data folder
	a.Log("Scanning data folder for JSON files...")
	jsonFiles, err := filepath.Glob(filepath.Join(gameInfo.DataPath, "*.json"))
	if err != nil {
		a.LogError("Failed to scan data folder")
		return err
	}
	a.Log(fmt.Sprintf("Found %d JSON files to patch", len(jsonFiles)))

	for _, jsonFile := range jsonFiles {
		err = PatchDataFile(a.ctx, jsonFile, patchInfo)
		if err != nil {
			a.LogError("Error patching file: " + filepath.Base(jsonFile))
			return err
		}
	}

	// Patch plugins.js
	pluginsJsPath := filepath.Join(gameInfo.JsPath, "plugins.js")
	err = UpdatePluginsJs(a.ctx, pluginsJsPath, patchInfo.Config.PluginsToPatch, patchInfo.Dictionary)
	if err != nil {
		a.LogError("Failed to update plugins.js")
		return err
	}

	// Apply replace rules
	for _, pluginToPatch := range patchInfo.Config.PluginsToPatch {
		for _, replaceRule := range pluginToPatch.ReplaceRules {
			err = ApplyReplaceRule(a.ctx, gameInfo.JsPath, pluginToPatch.Plugin, replaceRule)
			if err != nil {
				a.LogError("Failed to apply plugin replace rule")
				return err
			}
		}
	}

	a.Log("Reading system information...")
	systemInfoData, err := os.ReadFile(filepath.Join(gameInfo.DataPath, "system.json"))
	if err != nil {
		a.LogError("Failed to read system.json")
		return err
	}

	var systemInfo System
	if err := json.Unmarshal(systemInfoData, &systemInfo); err != nil {
		a.LogError("Failed to parse system.json")
		return err
	}

	a.Log("Looking for main screen image...")
	mainScreenImageName := systemInfo.Title1Name
	pngPath := filepath.Join(gameInfo.GameDir, "img", "titles1", mainScreenImageName+".png")
	if _, err := os.Stat(pngPath); os.IsNotExist(err) {
		pngPath = filepath.Join(gameInfo.GameDir, "img", "titles1", mainScreenImageName+".rpgmvp")
	}
	if _, err := os.Stat(pngPath); os.IsNotExist(err) {
		pngPath = filepath.Join(gameInfo.GameDir, "img", "titles1", mainScreenImageName+".png_")
	}
	if _, err := os.Stat(pngPath); os.IsNotExist(err) {
		a.LogError("Main screen image not found")
		return errors.New("main screen image not found")
	}

	if patchInfo.Config.CreditsLocation == "" {
		patchInfo.Config.CreditsLocation = "bottom_left"
	}

	a.Log("Adding credits to main screen image...")
	err = AddCreditsToResource(pngPath, systemInfo.EncryptionKey, patchInfo.Config.CreditsLocation)
	if err != nil {
		a.LogError("Failed to add credits")
		return err
	}

	a.LogSuccess("✓ Patch applied successfully!")
	return nil
}
