package service

import (
	"context"
	"errors"
	"htpatcher/internal/domain"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// StorageRepository interface for persistent storage
type StorageRepository interface {
	Load() (*domain.PersistentData, error)
	Save(data *domain.PersistentData) error
	GetDataPath() (string, error)
	DeleteData() error
}

// CollectionService handles game collection operations
type CollectionService struct {
	storage StorageRepository
	data    *domain.PersistentData
}

// NewCollectionService creates a new collection service
func NewCollectionService(storage StorageRepository) (*CollectionService, error) {
	data, err := storage.Load()
	if err != nil {
		return nil, err
	}
	return &CollectionService{
		storage: storage,
		data:    data,
	}, nil
}

// PrepareGameToAddToCollection opens a dialog to select a game to add
func (s *CollectionService) PrepareGameToAddToCollection(ctx context.Context) (*domain.LocatedGame, error) {
	exePath, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select the game executable",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Game Executable",
				Pattern:     "*.exe",
			},
		},
	})

	if err != nil || exePath == "" {
		return nil, errors.New("game selection aborted")
	}

	gameDir := filepath.Dir(exePath)
	locatedGame := &domain.LocatedGame{
		GameDir: gameDir,
		ExePath: exePath,
	}

	return locatedGame, nil
}

// AddGameToCollection adds a game to the user's collection
func (s *CollectionService) AddGameToCollection(locatedGame *domain.LocatedGame, rjCode string, friendlyName string, tags []string) error {
	locatedGame.RJCode = rjCode
	locatedGame.FriendlyName = friendlyName
	locatedGame.Tags = tags
	locatedGame.Id = uuid.NewString()
	s.data.LocatedGames = append(s.data.LocatedGames, *locatedGame)
	return s.storage.Save(s.data)
}

// GetDataPath returns the path to the persistent data file
func (s *CollectionService) GetDataPath() (string, error) {
	return s.storage.GetDataPath()
}

// DeleteData deletes the persistent data file
func (s *CollectionService) DeleteData() error {
	return s.storage.DeleteData()
}

// GetGamesCollection returns all games in the collection
func (s *CollectionService) GetGamesCollection() ([]domain.LocatedGame, error) {
	return s.data.LocatedGames, nil
}

// SetGameTranslated marks a game as translated or not
func (s *CollectionService) SetGameTranslated(id string, isTranslated bool) error {
	for i, game := range s.data.LocatedGames {
		if game.Id == id {
			s.data.LocatedGames[i].Translated = isTranslated
			return s.storage.Save(s.data)
		}
	}
	return errors.New("game not found")
}

// RemoveGameFromCollection removes a game from the collection
func (s *CollectionService) RemoveGameFromCollection(id string) error {
	for i, game := range s.data.LocatedGames {
		if game.Id == id {
			s.data.LocatedGames = append(s.data.LocatedGames[:i], s.data.LocatedGames[i+1:]...)
			break
		}
	}
	return s.storage.Save(s.data)
}

// UpdateGameMetadata updates the friendly name and tags for a game
func (s *CollectionService) UpdateGameMetadata(id string, friendlyName string, tags []string) error {
	for i, game := range s.data.LocatedGames {
		if game.Id == id {
			s.data.LocatedGames[i].FriendlyName = friendlyName
			s.data.LocatedGames[i].Tags = tags
			return s.storage.Save(s.data)
		}
	}
	return errors.New("game not found")
}

// GetGamesPerRow returns the games per row setting
func (s *CollectionService) GetGamesPerRow() int {
	if s.data.GamesPerRow == 0 {
		return 3
	}
	return s.data.GamesPerRow
}

// SetGamesPerRow sets the games per row setting
func (s *CollectionService) SetGamesPerRow(count int) error {
	if count != 3 && count != 4 {
		return errors.New("games per row must be 3 or 4")
	}
	s.data.GamesPerRow = count
	return s.storage.Save(s.data)
}

// SetGamePinned sets the pinned status of a game
func (s *CollectionService) SetGamePinned(id string, pinned bool) error {
	for i, game := range s.data.LocatedGames {
		if game.Id == id {
			s.data.LocatedGames[i].Pinned = pinned
			return s.storage.Save(s.data)
		}
	}
	return errors.New("game not found")
}

// SetGamePlayStatus sets the play status of a game
func (s *CollectionService) SetGamePlayStatus(id string, status string) error {
	validStatuses := map[string]bool{
		"unplayed": true,
		"playing":  true,
		"on-hold":  true,
		"finished": true,
		"given-up": true,
	}
	if !validStatuses[status] {
		return errors.New("invalid play status")
	}
	for i, game := range s.data.LocatedGames {
		if game.Id == id {
			s.data.LocatedGames[i].PlayStatus = status
			return s.storage.Save(s.data)
		}
	}
	return errors.New("game not found")
}
