package service

import (
	"encoding/json"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/util"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type UpdateService struct {
	logger Logger
}

func NewUpdateService(logger Logger) *UpdateService {
	return &UpdateService{logger: logger}
}

func (s *UpdateService) GetLatestReleaseInfo() (*domain.ReleaseInfo, error) {
	response, err := http.Get("https://api.github.com/repos/htranslations/htpatcher/releases/latest")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var releaseInfo domain.ReleaseInfo
	if err := json.Unmarshal(body, &releaseInfo); err != nil {
		return nil, err
	}
	return &releaseInfo, nil
}

func (s *UpdateService) GetNewReleaseInfo() (*domain.ReleaseInfo, error) {
	latestReleaseInfo, err := s.GetLatestReleaseInfo()
	if err != nil {
		return nil, err
	}

	latestVersion := strings.TrimPrefix(latestReleaseInfo.TagName, "v")
	latestVersionInt, err := strconv.Atoi(latestVersion)
	if err != nil {
		return nil, err
	}
	if latestVersionInt > Version {
		return latestReleaseInfo, nil
	}
	return nil, nil
}

func (s *UpdateService) GetCurrentVersion() int {
	return Version
}

// DownloadUpdate downloads the update executable to cache
func (s *UpdateService) DownloadUpdate(releaseInfo *domain.ReleaseInfo) error {
	// Find .exe asset
	asset := s.findExeAsset(releaseInfo.Assets)
	if asset == nil {
		return fmt.Errorf("no executable found in release assets")
	}

	// Get cache path
	updateExePath, err := util.GetUpdateExePath()
	if err != nil {
		return fmt.Errorf("failed to get cache path: %w", err)
	}

	s.logger.Info(fmt.Sprintf("Downloading update from %s", asset.BrowserDownloadURL))

	// Download with progress tracking
	resp, err := http.Get(asset.BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Create output file
	out, err := os.Create(updateExePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Download with progress tracking
	totalSize := resp.ContentLength
	downloaded := int64(0)
	lastEmit := time.Now()

	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := out.Write(buffer[:n]); writeErr != nil {
				os.Remove(updateExePath) // Clean up partial download
				return fmt.Errorf("write failed: %w", writeErr)
			}
			downloaded += int64(n)

			// Emit progress every 0.5 seconds
			if time.Since(lastEmit) >= 500*time.Millisecond {
				percentage := float64(downloaded) / float64(totalSize) * 100
				s.logger.Info(fmt.Sprintf("UPDATE_PROGRESS:%d:%d:%.1f", downloaded, totalSize, percentage))
				lastEmit = time.Now()
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Remove(updateExePath) // Clean up partial download
			return fmt.Errorf("download interrupted: %w", err)
		}
	}

	// Emit final progress
	s.logger.Info(fmt.Sprintf("UPDATE_PROGRESS:%d:%d:100.0", totalSize, totalSize))
	s.logger.Success("Download complete")
	return nil
}

// findExeAsset finds the .exe asset in release assets
func (s *UpdateService) findExeAsset(assets []domain.Asset) *domain.Asset {
	// First, look for exact name match
	for i := range assets {
		if strings.ToLower(assets[i].Name) == "htpatcher.exe" {
			return &assets[i]
		}
	}

	// Fallback: find any .exe file
	for i := range assets {
		if strings.HasSuffix(strings.ToLower(assets[i].Name), ".exe") {
			return &assets[i]
		}
	}

	return nil
}

// ApplyUpdate launches the downloaded update and exits the current app
func (s *UpdateService) ApplyUpdate() error {
	// Get paths
	updateExePath, err := util.GetUpdateExePath()
	if err != nil {
		return fmt.Errorf("failed to get update exe path: %w", err)
	}

	// Verify update exists
	if _, err := os.Stat(updateExePath); os.IsNotExist(err) {
		return fmt.Errorf("update file not found")
	}

	// Get current executable path
	currentExePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current exe path: %w", err)
	}

	s.logger.Info("Launching updater...")

	// Launch updater with --update flag
	cmd := exec.Command(updateExePath, "--update", currentExePath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch updater: %w", err)
	}

	// Exit immediately - updater will take over
	os.Exit(0)

	return nil // Never reached
}
