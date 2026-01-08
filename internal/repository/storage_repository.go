package repository

import (
	"encoding/json"
	"htpatcher/internal/domain"
	"os"
	"path/filepath"
)

// StorageRepository handles persistent data storage
type StorageRepository struct{}

// NewStorageRepository creates a new storage repository
func NewStorageRepository() *StorageRepository {
	return &StorageRepository{}
}

// Load loads persistent data from disk
func (r *StorageRepository) Load() (*domain.PersistentData, error) {
	path, err := r.getDataPath()
	if err != nil {
		return nil, err
	}

	// If the file does not exist, return empty data
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &domain.PersistentData{
			LocatedGames: []domain.LocatedGame{},
			GamesPerRow:  3,
		}, nil
	}

	// If the file exists, read it
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var persistentData domain.PersistentData
	if err := json.Unmarshal(data, &persistentData); err != nil {
		return nil, err
	}
	
	// Set default value for GamesPerRow if not present (backward compatibility)
	if persistentData.GamesPerRow == 0 {
		persistentData.GamesPerRow = 3
	}
	
	return &persistentData, nil
}

// Save saves persistent data to disk
func (r *StorageRepository) Save(persistentData *domain.PersistentData) error {
	path, err := r.getDataPath()
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

// GetDataPath returns the path to the persistent data file
func (r *StorageRepository) GetDataPath() (string, error) {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appDataDir, "htpatcher", "data.json"), nil
}

// getDataPath returns the path to the persistent data file (internal use)
func (r *StorageRepository) getDataPath() (string, error) {
	return r.GetDataPath()
}

// DeleteData deletes the persistent data file
func (r *StorageRepository) DeleteData() error {
	path, err := r.getDataPath()
	if err != nil {
		return err
	}
	
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // Already deleted, no error
	}
	
	return os.Remove(path)
}

