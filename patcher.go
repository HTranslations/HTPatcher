package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func getDataFileTypeMap(filePath string) string {
	fileTypeMap := map[string]string{
		"commonevents": "commonevents",
		"mapinfos":     "mapinfos",
		"actors":       "actors",
		"animations":   "animations",
		"armors":       "armors",
		"classes":      "classes",
		"enemies":      "enemies",
		"items":        "items",
		"skills":       "skills",
		"states":       "states",
		"system":       "system",
		"tilesets":     "tilesets",
		"troops":       "troops",
		"weapons":      "weapons",
	}

	filename := filepath.Base(filePath)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	filename = strings.ToLower(filename)

	if fileType, ok := fileTypeMap[filename]; ok {
		return fileType
	}
	if strings.HasPrefix(filename, "map") {
		return "map"
	}
	return filename
}

func PatchDataFile(ctx context.Context, filePath string, patchInfo PatchInfo) error {
	filename := filepath.Base(filePath)
	runtime.EventsEmit(ctx, "log", LogMessage{
		Message: "Patching: " + filename,
		Type:    "info",
	})
	fmt.Println("Patching file: ", filePath)
	fileType := getDataFileTypeMap(filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var patchedData []byte
	var patchError error

	switch fileType {
	case "actors":
		patchedData, patchError = PatchActors(data, patchInfo)
	case "armors":
		patchedData, patchError = PatchArmors(data, patchInfo)
	case "classes":
		patchedData, patchError = PatchClasses(data, patchInfo)
	case "commonevents":
		patchedData, patchError = PatchCommonEvents(data, patchInfo)
	case "enemies":
		patchedData, patchError = PatchEnemies(data, patchInfo)
	case "items":
		patchedData, patchError = PatchItems(data, patchInfo)
	case "map":
		patchedData, patchError = PatchMap(data, patchInfo)
	case "skills":
		patchedData, patchError = PatchSkills(data, patchInfo)
	case "states":
		patchedData, patchError = PatchStates(data, patchInfo)
	case "system":
		patchedData, patchError = PatchSystem(data, patchInfo)
	case "troops":
		patchedData, patchError = PatchTroops(data, patchInfo)
	case "weapons":
		patchedData, patchError = PatchWeapons(data, patchInfo)
	default:
		return nil
	}

	if patchError != nil {
		return patchError
	}

	return os.WriteFile(filePath, patchedData, 0644)
}
