package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type PersistentData struct {
	LocatedGames []LocatedGame `json:"locatedGames"`
}

type LocatedGame struct {
	Id         string `json:"id"`
	GameDir    string `json:"gameDir"`
	ExePath    string `json:"exePath"`
	RJCode     string `json:"rjCode"`
	Translated bool   `json:"translated"`
}

func GetPersistentDataPath() (string, error) {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appDataDir, "htpatcher", "data.json"), nil
}

func GetPersistentData() (*PersistentData, error) {
	path, err := GetPersistentDataPath()
	if err != nil {
		return nil, err
	}

	// If the file does not exist, return an empty PersistentData
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &PersistentData{
			LocatedGames: []LocatedGame{},
		}, nil
	}

	// If the file exists, read it and return the PersistentData
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var persistentData PersistentData
	if err := json.Unmarshal(data, &persistentData); err != nil {
		return nil, err
	}
	return &persistentData, nil
}

func SavePersistentData(persistentData *PersistentData) error {
	path, err := GetPersistentDataPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(persistentData, "", "  ")
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(path), 0755)

	return os.WriteFile(path, data, 0644)
}
