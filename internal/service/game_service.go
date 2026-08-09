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

	// Check for VX/VX Ace games first.
	rgss3aPath := filepath.Join(gameInfo.GameDir, "Game.rgss3a")
	rgss2aPath := filepath.Join(gameInfo.GameDir, "Game.rgss2a")
	dataPathRGSS := filepath.Join(gameInfo.GameDir, "Data")

	if _, err := os.Stat(rgss3aPath); err == nil {
		return s.getRGSSGameInfo(&gameInfo, "vxace", rgss3aPath, dataPathRGSS)
	}
	if _, err := os.Stat(rgss2aPath); err == nil {
		return s.getRGSSGameInfo(&gameInfo, "vx", rgss2aPath, dataPathRGSS)
	}
	if _, err := os.Stat(dataPathRGSS); err == nil {
		files, err := os.ReadDir(dataPathRGSS)
		if err == nil {
			for _, f := range files {
				if strings.HasSuffix(strings.ToLower(f.Name()), ".rvdata2") {
					return s.getRGSSGameInfo(&gameInfo, "vxace", rgss3aPath, dataPathRGSS)
				}
				if strings.HasSuffix(strings.ToLower(f.Name()), ".rvdata") {
					return s.getRGSSGameInfo(&gameInfo, "vx", rgss2aPath, dataPathRGSS)
				}
			}
		}
	}

	// Standard MV/MZ game detection
	return s.getMVMZGameInfo(&gameInfo)
}

// getRGSSGameInfo handles VX and VX Ace game detection.
func (s *GameService) getRGSSGameInfo(gameInfo *domain.GameInfo, gameVersion, archivePath, dataPath string) (*domain.GameInfo, error) {
	gameInfo.GameVersion = gameVersion
	gameInfo.DataPath = dataPath
	engineName := "VX Ace"
	dataFile := "System.rvdata2"
	archiveName := "Game.rgss3a"
	if gameVersion == "vx" {
		engineName = "VX"
		dataFile = "System.rvdata"
		archiveName = "Game.rgss2a"
	}

	// VX Ace doesn't have js or img paths in the same way
	gameInfo.JsPath = ""
	gameInfo.ImgPath = filepath.Join(gameInfo.GameDir, "Graphics")

	systemPath := filepath.Join(dataPath, dataFile)

	// If archive exists and System.rvdata2 doesn't exist in Data folder, read from archive
	if _, err := os.Stat(systemPath); os.IsNotExist(err) {
		if _, err := os.Stat(archivePath); err == nil {
			// Need to read from archive
			s.logger.Info("Reading system data from " + archiveName + " archive")
			title, err := s.readTitleFromArchive(archivePath, dataFile)
			if err != nil {
				s.logger.Warn(fmt.Sprintf("Failed to read title from archive: %v", err))
				gameInfo.GameTitle = "Unknown " + engineName + " Game"
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
			s.logger.Warn(fmt.Sprintf("Failed to read %s: %v", dataFile, err))
			gameInfo.GameTitle = "Unknown " + engineName + " Game"
			return gameInfo, nil
		}

		raw, err := marshal.Parse(data)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("Failed to parse %s: %v", dataFile, err))
			gameInfo.GameTitle = "Unknown " + engineName + " Game"
			return gameInfo, nil
		}

		gameInfo.GameTitle = rpgmakervxa.GetSystemTitle(raw)
		if gameInfo.GameTitle == "" {
			gameInfo.GameTitle = "Unknown " + engineName + " Game"
		}
	} else {
		gameInfo.GameTitle = "Unknown " + engineName + " Game"
	}

	s.logger.Info(fmt.Sprintf("%s game detected: \"%s\"", engineName, gameInfo.GameTitle))
	return gameInfo, nil
}

// readTitleFromArchive reads the game title from an RGSSAD archive.
func (s *GameService) readTitleFromArchive(archivePath, dataFile string) (string, error) {
	archive, err := rgss3a.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer archive.Close()

	data, err := archive.ReadFile("Data/" + dataFile)
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

	gameInfo.DataCipher = util.DetectDataCipher(jsPath)
	if gameInfo.DataCipher != nil {
		s.logger.Info(fmt.Sprintf("Detected data file XOR obfuscation (key=%d)", gameInfo.DataCipher.K))
	}

	gameInfo.WrapperCipher = util.DetectWrapperCipher(dataPath)
	if gameInfo.WrapperCipher != nil {
		s.logger.Info("Detected custom data file wrapper encryption")
	}

	systemInfoData, err := util.ReadDataFile(filepath.Join(gameInfo.DataPath, "System.json"), gameInfo.DataCipher, gameInfo.WrapperCipher)
	if err != nil {
		s.logger.Error("Failed to read System.json")
		return nil, err
	}

	var systemInfo rpgmaker.System
	if err := json.Unmarshal(systemInfoData, &systemInfo); err != nil {
		s.logger.Error("Failed to parse System.json")
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
