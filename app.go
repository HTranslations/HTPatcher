package main

import (
	"context"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/repository"
	"htpatcher/internal/service"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx               context.Context
	gameService       *service.GameService
	patchService      *service.PatchService
	backupService     *service.BackupService
	collectionService *service.CollectionService
	downloadService   *service.DownloadService
	updateService     *service.UpdateService
	exportService     *service.ExportService
	justUpdated       bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.WindowMaximise(a.ctx)

	// Emit update success event if just updated
	if a.justUpdated {
		runtime.EventsEmit(a.ctx, "app:updated", service.Version)
	}

	// Initialize logger
	logger := &AppLogger{ctx: ctx}

	// Initialize repositories
	patchRepo := repository.NewPatchRepository()
	storageRepo := repository.NewStorageRepository()

	// Initialize services
	a.gameService = service.NewGameService(logger)
	a.patchService = service.NewPatchService(patchRepo, logger)
	a.backupService = service.NewBackupService(logger)
	a.downloadService = service.NewDownloadService(patchRepo, logger)
	a.updateService = service.NewUpdateService(logger)
	a.exportService = service.NewExportService(logger)

	collectionService, err := service.NewCollectionService(storageRepo)
	if err != nil {
		a.LogError("Failed to initialize collection service")
		return
	}
	a.collectionService = collectionService
}

// AppLogger implements the Logger interface for services
type AppLogger struct {
	ctx context.Context
}

func (l *AppLogger) Info(message string) {
	runtime.EventsEmit(l.ctx, "log", LogMessage{
		Message: message,
		Type:    "info",
	})
}

func (l *AppLogger) Success(message string) {
	runtime.EventsEmit(l.ctx, "log", LogMessage{
		Message: message,
		Type:    "success",
	})
}

func (l *AppLogger) Error(message string) {
	runtime.EventsEmit(l.ctx, "log", LogMessage{
		Message: message,
		Type:    "error",
	})
}

// LogMessage represents a log message sent to the frontend
type LogMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"` // "info", "success", "error"
}

// Log logs an info message
func (a *App) Log(message string) {
	runtime.EventsEmit(a.ctx, "log", LogMessage{
		Message: message,
		Type:    "info",
	})
}

// LogSuccess logs a success message
func (a *App) LogSuccess(message string) {
	runtime.EventsEmit(a.ctx, "log", LogMessage{
		Message: message,
		Type:    "success",
	})
}

// LogError logs an error message
func (a *App) LogError(message string) {
	runtime.EventsEmit(a.ctx, "log", LogMessage{
		Message: message,
		Type:    "error",
	})
}

// ===== Game Service Methods =====

// SelectGameExeFile opens a dialog to select a game executable
func (a *App) SelectGameExeFile() (*domain.GameInfo, error) {
	return a.gameService.SelectGameExeFile(a.ctx)
}

// GetGameInfoFromExePath gets game info from an executable path
func (a *App) GetGameInfoFromExePath(exePath string) (*domain.GameInfo, error) {
	return a.gameService.GetGameInfoFromExePath(exePath)
}

// LaunchGameFromPath launches a game from its path
func (a *App) LaunchGameFromPath(exePath string) error {
	return a.gameService.LaunchGame(exePath)
}

// OpenFolder opens a folder in the file explorer
func (a *App) OpenFolder(folderPath string) error {
	return a.gameService.OpenFolder(folderPath)
}

// ===== Patch Service Methods =====

// SelectPatchFile opens a dialog to select a patch file
func (a *App) SelectPatchFile() (*domain.PatchInfo, error) {
	return a.patchService.SelectPatchFile(a.ctx)
}

// FetchAllPatches fetches all available patches
func (a *App) FetchAllPatches() ([]domain.PatchEntry, error) {
	return a.patchService.FetchAllPatches()
}

// ApplyPatch applies a patch to a game
func (a *App) ApplyPatch(gameInfo domain.GameInfo, patchInfo domain.PatchInfo, launchAfterPatch bool, backupBeforePatch bool) error {
	if backupBeforePatch {
		a.Log("Backing up game data...")
		err := a.backupService.BackupGameData(&gameInfo, &patchInfo)
		if err != nil {
			a.LogError("Failed to backup game data")
			return err
		}
	}

	err := a.patchService.ApplyPatch(a.ctx, &gameInfo, &patchInfo)
	if err != nil {
		return err
	}

	if launchAfterPatch {
		a.Log("Launching game...")
		err = a.gameService.LaunchGame(gameInfo.ExePath)
		if err != nil {
			a.LogError("Failed to launch game")
			return err
		}
		a.Log("Game launched successfully!")
	}

	return nil
}

// ===== Download Service Methods =====

// DownloadPatch downloads a patch
func (a *App) DownloadPatch(patchDownloadId string) (*domain.PatchInfo, error) {
	return a.downloadService.DownloadPatch(patchDownloadId, func(filePath string) (*domain.PatchInfo, error) {
		patchInfo, err := a.patchService.LoadPatchInfo(filePath)
		if err != nil {
			return nil, err
		}
		// Clean up temp file
		os.Remove(filePath)
		a.Log(fmt.Sprintf("Removed %s", filePath))
		return patchInfo, nil
	})
}

// ===== Backup Service Methods =====

// RestoreGameBackup restores a game backup
func (a *App) RestoreGameBackup(gameInfo domain.GameInfo) error {
	a.Log("Starting backup restoration...")
	err := a.backupService.RestoreBackup(&gameInfo)
	if err != nil {
		a.LogError("Failed to restore backup")
		return err
	}
	a.LogSuccess("✓ Backup restored successfully!")
	return nil
}

// ===== Collection Service Methods =====

// PrepareGameToAddToCollection prepares a game to be added to the collection
func (a *App) PrepareGameToAddToCollection() (*domain.LocatedGame, error) {
	return a.collectionService.PrepareGameToAddToCollection(a.ctx)
}

// AddGameToCollection adds a game to the collection
func (a *App) AddGameToCollection(locatedGame *domain.LocatedGame, rjCode string, friendlyName string, tags []string) error {
	return a.collectionService.AddGameToCollection(locatedGame, rjCode, friendlyName, tags)
}

// GetGamesCollection returns all games in the collection
func (a *App) GetGamesCollection() ([]domain.LocatedGame, error) {
	return a.collectionService.GetGamesCollection()
}

// SetGameTranslated marks a game as translated
func (a *App) SetGameTranslated(id string, isTranslated bool) error {
	return a.collectionService.SetGameTranslated(id, isTranslated)
}

// RemoveGameFromCollection removes a game from the collection
func (a *App) RemoveGameFromCollection(id string) error {
	return a.collectionService.RemoveGameFromCollection(id)
}

// GetPersistentDataPath returns the path to the persistent data file
func (a *App) GetPersistentDataPath() (string, error) {
	return a.collectionService.GetDataPath()
}

// DeletePersistentData deletes the persistent data file
func (a *App) DeletePersistentData() error {
	err := a.collectionService.DeleteData()
	if err != nil {
		return err
	}
	// Reload the collection service with empty data
	a.collectionService, err = service.NewCollectionService(repository.NewStorageRepository())
	return err
}

// UpdateGameMetadata updates the friendly name and tags for a game
func (a *App) UpdateGameMetadata(gameId string, friendlyName string, tags []string) error {
	return a.collectionService.UpdateGameMetadata(gameId, friendlyName, tags)
}

// GetGamesPerRow returns the games per row setting
func (a *App) GetGamesPerRow() int {
	return a.collectionService.GetGamesPerRow()
}

// SetGamesPerRow sets the games per row setting
func (a *App) SetGamesPerRow(count int) error {
	return a.collectionService.SetGamesPerRow(count)
}

// SetGamePinned sets the pinned status of a game
func (a *App) SetGamePinned(id string, pinned bool) error {
	return a.collectionService.SetGamePinned(id, pinned)
}

// SetGamePlayStatus sets the play status of a game
func (a *App) SetGamePlayStatus(id string, status string) error {
	return a.collectionService.SetGamePlayStatus(id, status)
}

// ===== Update Service Methods =====

// CheckForUpdate checks if there's a new version available
func (a *App) CheckForUpdate() (*domain.ReleaseInfo, error) {
	return a.updateService.GetNewReleaseInfo()
}

// GetCurrentVersion returns the current application version
func (a *App) GetCurrentVersion() int {
	return a.updateService.GetCurrentVersion()
}

// GetLatestReleaseInfo returns the latest release information from GitHub
func (a *App) GetLatestReleaseInfo() (*domain.ReleaseInfo, error) {
	return a.updateService.GetLatestReleaseInfo()
}

// DownloadUpdate downloads the new version to cache
func (a *App) DownloadUpdate(releaseInfo *domain.ReleaseInfo) error {
	return a.updateService.DownloadUpdate(releaseInfo)
}

// ApplyUpdate applies the downloaded update and restarts the app
func (a *App) ApplyUpdate() error {
	return a.updateService.ApplyUpdate()
}

// ===== Export Service Methods =====

// ExportPatchedFiles exports patched files to a ZIP archive
func (a *App) ExportPatchedFiles(gameDir string, friendlyName string) error {
	return a.exportService.ExportPatchedFiles(a.ctx, gameDir, friendlyName)
}
