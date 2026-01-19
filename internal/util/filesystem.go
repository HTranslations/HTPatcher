package util

import (
	"os"
	"path/filepath"
	"strings"
)

func ListFilesWithExtensions(path string, extensions []string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	filesWithExtension := []string{}
	for _, file := range files {
		for _, extension := range extensions {
			if strings.HasSuffix(file.Name(), extension) {
				filesWithExtension = append(filesWithExtension, filepath.Join(path, file.Name()))
				break
			}
		}
	}
	return filesWithExtension, nil
}

func ListFilesWithExtension(path string, extension string) ([]string, error) {
	return ListFilesWithExtensions(path, []string{extension})
}
