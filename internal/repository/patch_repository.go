package repository

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"htpatcher/internal/domain"
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
	return *dictionary, nil
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

// Download downloads a patch from the download service
func (r *PatchRepository) Download(patchDownloadId string) (string, error) {
	url := fmt.Sprintf("https://cybersharing.net/api/containers/%s", patchDownloadId)

	response, err := http.Post(url, "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var container ContainerResponse
	if err := json.Unmarshal(body, &container); err != nil {
		return "", err
	}

	if len(container.Uploads) == 0 {
		return "", errors.New("no uploads found in container")
	}

	// Find the .htpatch file
	var patchUpload *Upload
	for _, upload := range container.Uploads {
		if strings.HasSuffix(upload.FileName, ".htpatch") {
			patchUpload = &upload
			break
		}
	}

	if patchUpload == nil {
		return "", errors.New("no patch upload found in container")
	}

	downloadUrl := fmt.Sprintf("https://cybersharing.net/api/download/file/%s/%s/%s/%s",
		container.ID, patchUpload.ID, container.Signature, patchUpload.FileName)

	downloadResponse, err := http.Get(downloadUrl)
	if err != nil {
		return "", err
	}
	defer downloadResponse.Body.Close()

	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, patchUpload.FileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, downloadResponse.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// Helper types for download API
type ContainerResponse struct {
	ID        string   `json:"id"`
	Signature string   `json:"signature"`
	Uploads   []Upload `json:"uploads"`
}

type Upload struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
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




