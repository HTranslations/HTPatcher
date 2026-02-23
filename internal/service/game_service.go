package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/domain/rpgmaker"
	"htpatcher/internal/domain/rpgmakervxa"
	"htpatcher/internal/marshal"
	"htpatcher/internal/rgss3a"
	"htpatcher/internal/util"
	"os"
	"path/filepath"
	"strings"

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

	// Check for VX Ace game first
	rgss3aPath := filepath.Join(gameInfo.GameDir, "Game.rgss3a")
	dataPathVXA := filepath.Join(gameInfo.GameDir, "Data")

	// Check if this is a VX Ace game
	isVXAce := false
	if _, err := os.Stat(rgss3aPath); err == nil {
		// Has Game.rgss3a - definitely VX Ace
		isVXAce = true
	} else if _, err := os.Stat(dataPathVXA); err == nil {
		// Check for .rvdata2 files in Data folder
		files, err := os.ReadDir(dataPathVXA)
		if err == nil {
			for _, f := range files {
				if strings.HasSuffix(strings.ToLower(f.Name()), ".rvdata2") {
					isVXAce = true
					break
				}
			}
		}
	}

	if isVXAce {
		return s.getVXAceGameInfo(&gameInfo, rgss3aPath, dataPathVXA)
	}

	// Standard MV/MZ game detection
	return s.getMVMZGameInfo(&gameInfo)
}

// getVXAceGameInfo handles VX Ace game detection
func (s *GameService) getVXAceGameInfo(gameInfo *domain.GameInfo, rgss3aPath, dataPath string) (*domain.GameInfo, error) {
	gameInfo.GameVersion = "vxace"
	gameInfo.DataPath = dataPath

	// VX Ace doesn't have js or img paths in the same way
	gameInfo.JsPath = ""
	gameInfo.ImgPath = filepath.Join(gameInfo.GameDir, "Graphics")

	// Try to read game title from System.rvdata2
	systemPath := filepath.Join(dataPath, "System.rvdata2")

	// If archive exists and System.rvdata2 doesn't exist in Data folder, read from archive
	if _, err := os.Stat(systemPath); os.IsNotExist(err) {
		if _, err := os.Stat(rgss3aPath); err == nil {
			// Need to read from archive
			s.logger.Info("Reading system data from Game.rgss3a archive")
			title, err := s.readTitleFromArchive(rgss3aPath)
			if err != nil {
				s.logger.Warn(fmt.Sprintf("Failed to read title from archive: %v", err))
				gameInfo.GameTitle = "Unknown VX Ace Game"
			} else {
				gameInfo.GameTitle = title
			}
			return gameInfo, nil
		}
	}

	// Read from extracted file
	if _, err := os.Stat(systemPath); err == nil {
		data, err := os.ReadFile(systemPath)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("Failed to read System.rvdata2: %v", err))
			gameInfo.GameTitle = "Unknown VX Ace Game"
			return gameInfo, nil
		}

		raw, err := marshal.Parse(data)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("Failed to parse System.rvdata2: %v", err))
			gameInfo.GameTitle = "Unknown VX Ace Game"
			return gameInfo, nil
		}

		gameInfo.GameTitle = rpgmakervxa.GetSystemTitle(raw)
		if gameInfo.GameTitle == "" {
			gameInfo.GameTitle = "Unknown VX Ace Game"
		}
	} else {
		gameInfo.GameTitle = "Unknown VX Ace Game"
	}

	s.logger.Info(fmt.Sprintf("VX Ace game detected: \"%s\"", gameInfo.GameTitle))
	return gameInfo, nil
}

// readTitleFromArchive reads the game title from a Game.rgss3a archive
func (s *GameService) readTitleFromArchive(archivePath string) (string, error) {
	archive, err := rgss3a.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer archive.Close()

	data, err := archive.ReadFile("Data/System.rvdata2")
	if err != nil {
		return "", err
	}

	raw, err := marshal.Parse(data)
	if err != nil {
		return "", err
	}

	return rpgmakervxa.GetSystemTitle(raw), nil
}

// getMVMZGameInfo handles MV/MZ game detection
func (s *GameService) getMVMZGameInfo(gameInfo *domain.GameInfo) (*domain.GameInfo, error) {
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
	// Detect MV vs MZ by checking for MZ-specific core file
	if _, err := os.Stat(filepath.Join(jsPath, "rmmz_core.js")); err == nil {
		gameInfo.GameVersion = "mz"
	} else {
		gameInfo.GameVersion = "mv"
	}

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

	return gameInfo, nil
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
