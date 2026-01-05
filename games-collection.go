package main

import (
	"errors"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) PrepareGameToAddToCollection() (*LocatedGame, error) {
	exePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select the game executable",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Game Executable",
				Pattern:     "*.exe",
			},
		},
	})

	if err != nil || exePath == "" {
		return nil, errors.New("Game selection aborted")
	}

	gameDir := filepath.Dir(exePath)

	locatedGame := &LocatedGame{
		GameDir: gameDir,
		ExePath: exePath,
	}

	return locatedGame, nil
}

func (a *App) AddGameToCollection(locatedGame *LocatedGame, rjCode string) error {
	locatedGame.RJCode = rjCode
	locatedGame.Id = uuid.NewString()
	a.persistentData.LocatedGames = append(a.persistentData.LocatedGames, *locatedGame)
	return SavePersistentData(a.persistentData)
}

func (a *App) GetGamesCollection() ([]LocatedGame, error) {
	return a.persistentData.LocatedGames, nil
}

func (a *App) SetGameTranslated(id string, isTranslated bool) error {
	for i, game := range a.persistentData.LocatedGames {
		if game.Id == id {
			a.persistentData.LocatedGames[i].Translated = isTranslated
			return SavePersistentData(a.persistentData)
		}
	}
	return errors.New("Game not found")
}

func (a *App) RemoveGameFromCollection(id string) error {
	for i, game := range a.persistentData.LocatedGames {
		if game.Id == id {
			a.persistentData.LocatedGames = append(a.persistentData.LocatedGames[:i], a.persistentData.LocatedGames[i+1:]...)
			break
		}
	}
	return SavePersistentData(a.persistentData)
}
