package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/domain/rpgmaker"
	"htpatcher/internal/util"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// GameService handles game-related operations
type GameService struct {
	logger Logger
}

// Logger interface for logging operations
type Logger interface {
	Info(message string)
	Success(message string)
	Error(message string)
	Warn(message string)
}

// NewGameService creates a new game service
func NewGameService(logger Logger) *GameService {
	return &GameService{logger: logger}
}

// SelectGameExeFile opens a file dialog to select a game executable
func (s *GameService) SelectGameExeFile(ctx context.Context) (*domain.GameInfo, error) {
	filePath, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
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

	return s.GetGameInfoFromExePath(filePath)
}

// GetGameInfoFromExePath extracts game information from an executable path
func (s *GameService) GetGameInfoFromExePath(filePath string) (*domain.GameInfo, error) {
	gameInfo := domain.GameInfo{}

	// Set game paths
	gameInfo.ExePath = filePath
	gameInfo.GameDir = filepath.Dir(filePath)
	s.logger.Info(fmt.Sprintf("Game directory: %s", gameInfo.GameDir))

	// Set data and js paths
	dataPath := filepath.Join(gameInfo.GameDir, "data")
	imgPath := filepath.Join(gameInfo.GameDir, "img")
	jsPath := filepath.Join(gameInfo.GameDir, "js")
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		dataPath = filepath.Join(gameInfo.GameDir, "www", "data")
		imgPath = filepath.Join(gameInfo.GameDir, "www", "img")
		jsPath = filepath.Join(gameInfo.GameDir, "www", "js")
	}

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		s.logger.Error("Data directory not found")
		return nil, errors.New("data directory not found")
	}
	if _, err := os.Stat(jsPath); os.IsNotExist(err) {
		s.logger.Error("JS directory not found")
		return nil, errors.New("js directory not found")
	}
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		s.logger.Error("IMG directory not found")
		return nil, errors.New("img directory not found")
	}

	gameInfo.DataPath = dataPath
	gameInfo.JsPath = jsPath
	gameInfo.ImgPath = imgPath

	systemInfoData, err := os.ReadFile(filepath.Join(gameInfo.DataPath, "system.json"))
	if err != nil {
		s.logger.Error("Failed to read system.json")
		return nil, err
	}

	var systemInfo rpgmaker.System
	if err := json.Unmarshal(systemInfoData, &systemInfo); err != nil {
		s.logger.Error("Failed to parse system.json")
		return nil, err
	}

	gameInfo.GameTitle = systemInfo.GameTitle
	s.logger.Info(fmt.Sprintf("Game title: \"%s\"", gameInfo.GameTitle))

	return &gameInfo, nil
}

// LaunchGame launches a game executable
func (s *GameService) LaunchGame(exePath string) error {
	workingDir := filepath.Dir(exePath)
	return util.LaunchExecutable(exePath, workingDir)
}

// OpenFolder opens a folder in the file explorer
func (s *GameService) OpenFolder(folderPath string) error {
	return util.OpenFolder(folderPath)
}
