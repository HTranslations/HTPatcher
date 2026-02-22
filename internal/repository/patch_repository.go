package repository

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// PatchRepository handles patch file operations
type PatchRepository struct{}

// NewPatchRepository creates a new patch repository
func NewPatchRepository() *PatchRepository {
	return &PatchRepository{}
}

// Open opens a patch file and returns a zip reader
func (r *PatchRepository) Open(path string) (*zip.ReadCloser, error) {
	return zip.OpenReader(path)
}

// ReadDictionary reads the translation dictionary from a patch
func (r *PatchRepository) ReadDictionary(zipReader *zip.ReadCloser) (map[string]string, error) {
	dictionary, err := readJSONFromZip[map[string]string](zipReader, "dictionary.json")
	if err != nil {
		return nil, err
	}

	// Normalize dictionary keys to match GetTranslationKey output
	normalized := make(map[string]string, len(*dictionary))
	for key, value := range *dictionary {
		normalized[util.GetTranslationKey(key)] = value
	}
	return normalized, nil
}

// ReadConfig reads the patch configuration
func (r *PatchRepository) ReadConfig(zipReader *zip.ReadCloser) (*domain.Config, error) {
	return readJSONFromZip[domain.Config](zipReader, "config.json")
}

// GetAllOverrides lists all override files in the patch
func (r *PatchRepository) GetAllOverrides(zipReader *zip.ReadCloser) ([]string, error) {
	overrides := []string{}
	for _, f := range zipReader.File {
		if strings.HasPrefix(f.Name, "overrides/") && f.Mode().IsRegular() {
			overrides = append(overrides, strings.TrimPrefix(f.Name, "overrides/"))
		}
	}
	return overrides, nil
}

// ReadFileFromZip reads a specific file from the patch
func (r *PatchRepository) ReadFileFromZip(zipReader *zip.ReadCloser, path string) ([]byte, error) {
	path = strings.ReplaceAll(path, "\\", "/")
	for _, f := range zipReader.File {
		if f.Name == path {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, errors.New("file " + path + " not found")
}

// Download downloads a patch file from a direct URL
func (r *PatchRepository) Download(url string, fileName string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", response.StatusCode)
	}

	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// readJSONFromZip is a generic helper to read and unmarshal JSON from a zip file
func readJSONFromZip[T any](zipReader *zip.ReadCloser, name string) (*T, error) {
	for _, f := range zipReader.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			b, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			var v T
			return &v, json.Unmarshal(b, &v)
		}
	}
	return nil, errors.New(name + " not found")
}




