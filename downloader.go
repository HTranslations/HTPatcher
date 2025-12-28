package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ContainerResponse struct {
	ID        string   `json:"id"`
	Signature string   `json:"signature"`
	Uploads   []Upload `json:"uploads"`
}

type Upload struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
}

func DownloadPatch(patchDownloadId string) (string, error) {
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

	// upload is any file that ends with .htpatch
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

	downloadUrl := fmt.Sprintf("https://cybersharing.net/api/download/file/%s/%s/%s/%s", container.ID, patchUpload.ID, container.Signature, patchUpload.FileName)

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
