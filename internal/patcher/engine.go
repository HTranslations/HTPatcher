package patcher

import (
	"context"
	"htpatcher/internal/domain"
	"htpatcher/internal/domain/rpgmaker"
	"htpatcher/internal/util"
	"path/filepath"
	"strings"
)

// Engine handles all patching operations
type Engine struct {
	logger Logger
}

// Logger interface for logging operations
type Logger interface {
	Info(message string)
	Success(message string)
	Error(message string)
	Warn(message string)
}

// NewEngine creates a new patcher engine
func NewEngine(logger Logger) *Engine {
	return &Engine{
		logger: logger,
	}
}

// PatchDataFile patches a single data file based on its type. If cipher is
// non-nil, the file is XOR-decrypted before patching and re-encrypted on write.
// If wrapperCipher is non-nil, it handles the {uid,bid,data} wrapper format.
func (e *Engine) PatchDataFile(ctx context.Context, filePath string, patchInfo *domain.PatchInfo, cipher *util.DataCipher, wrapperCipher *util.WrapperCipher) error {
	filename := filepath.Base(filePath)
	e.logger.Info("Patching: " + filename)

	fileType := getDataFileTypeMap(filePath)
	data, err := util.ReadDataFile(filePath, cipher, wrapperCipher)
	if err != nil {
		return err
	}

	var patchedData []byte
	var patchError error

	switch fileType {
	case "actors":
		patchedData, patchError = patchActors(data, patchInfo)
	case "armors":
		patchedData, patchError = patchArmors(data, patchInfo)
	case "classes":
		patchedData, patchError = patchClasses(data, patchInfo)
	case "commonevents":
		patchedData, patchError = patchCommonEvents(data, patchInfo)
	case "enemies":
		patchedData, patchError = patchEnemies(data, patchInfo)
	case "items":
		patchedData, patchError = patchItems(data, patchInfo)
	case "map":
		patchedData, patchError = patchMap(data, patchInfo)
	case "mapinfos":
		patchedData, patchError = patchMapInfos(data, patchInfo)
	case "skills":
		patchedData, patchError = patchSkills(data, patchInfo)
	case "states":
		patchedData, patchError = patchStates(data, patchInfo)
	case "system":
		patchedData, patchError = patchSystem(data, patchInfo)
	case "troops":
		patchedData, patchError = patchTroops(data, patchInfo)
	case "weapons":
		patchedData, patchError = patchWeapons(data, patchInfo)
	default:
		return nil
	}

	if patchError != nil {
		return patchError
	}

	return util.WriteDataFile(filePath, patchedData, cipher, wrapperCipher)
}

// getDataFileTypeMap determines the file type from the filename
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

// PatchCommands patches event commands (used by maps, common events, troops)
func (e *Engine) PatchCommands(commands []*rpgmaker.EventCommand, patchInfo *domain.PatchInfo) ([]*rpgmaker.EventCommand, error) {
	return patchCommands(commands, patchInfo)
}
