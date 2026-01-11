package patcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/util"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// PluginPatcher handles plugin patching operations
type PluginPatcher struct {
	logger Logger
}

// NewPluginPatcher creates a new plugin patcher
func NewPluginPatcher(logger Logger) *PluginPatcher {
	return &PluginPatcher{logger: logger}
}

// ApplyReplaceRule applies a replace rule to a plugin file
func (p *PluginPatcher) ApplyReplaceRule(ctx context.Context, jsPath string, pluginName string, replaceRule domain.PluginReplaceRule, ruleIndex int) error {
	p.logger.Info("Applying replace rule on plugin " + pluginName)

	pluginPath := filepath.Join(jsPath, "plugins", pluginName+".js")
	data, err := os.ReadFile(pluginPath)
	if err != nil {
		return err
	}

	// Normalize line endings by removing \r for cross-platform compatibility
	normalizedData := bytes.ReplaceAll(data, []byte("\r"), []byte(""))
	normalizedMatch := bytes.ReplaceAll([]byte(replaceRule.Match), []byte("\r"), []byte(""))
	normalizedReplace := bytes.ReplaceAll([]byte(replaceRule.Replace), []byte("\r"), []byte(""))

	// Check if match string exists in the normalized file
	if !bytes.Contains(normalizedData, normalizedMatch) {
		p.logger.Warn(fmt.Sprintf("Replace rule #%d was not applied on plugin %s", ruleIndex, pluginName))
	}

	patchedData := bytes.ReplaceAll(normalizedData, normalizedMatch, normalizedReplace)
	return os.WriteFile(pluginPath, patchedData, 0644)
}

// UpdatePluginsJs updates the plugins.js file with translated plugin parameters
func (p *PluginPatcher) UpdatePluginsJs(ctx context.Context, pluginsJsPath string, pluginsToPatch []domain.PluginToPatch, dictionary map[string]string) error {
	data, err := os.ReadFile(pluginsJsPath)
	if err != nil {
		return err
	}

	jsContent := string(data)
	startIndex := strings.Index(jsContent, "[")
	endIndex := strings.LastIndex(jsContent, "]")
	if startIndex == -1 || endIndex == -1 {
		return nil
	}

	pluginsJson := jsContent[startIndex : endIndex+1]
	var plugins []domain.PluginData
	err = json.Unmarshal([]byte(pluginsJson), &plugins)
	if err != nil {
		return err
	}

	p.logger.Info("Updating plugins data")

	for i := range plugins {
		for _, pluginToPatch := range pluginsToPatch {
			if plugins[i].Name == pluginToPatch.Plugin && pluginToPatch.ParametersPatchScript != "" {
				p.logger.Info("Patching plugin data of: " + plugins[i].Name)

				L := lua.NewState()
				defer L.Close()
				L.SetGlobal("getTranslationByKey", L.NewFunction(makeGetTranslationByKey(dictionary)))
				L.SetGlobal("jsonDecode", L.NewFunction(jsonDecode))
				L.SetGlobal("jsonEncode", L.NewFunction(jsonEncode))
				if err := L.DoString(pluginToPatch.ParametersPatchScript); err != nil {
					return err
				}

				fn := L.GetGlobal("patch")
				L.Push(fn)
				L.Push(lua.LString(string(plugins[i].Parameters)))
				if err := L.PCall(1, 1, nil); err != nil {
					return err
				}

				patchedParams := L.ToString(-1)
				L.Pop(1)
				plugins[i].Parameters = json.RawMessage(patchedParams)
			}
		}
	}

	patchedPluginsJson, err := json.Marshal(plugins)
	if err != nil {
		return err
	}

	before := jsContent[:startIndex]
	after := jsContent[endIndex+1:]
	patchedData := before + string(patchedPluginsJson) + after

	return os.WriteFile(pluginsJsPath, []byte(patchedData), 0644)
}

// Lua helper functions
func makeGetTranslationByKey(dictionary map[string]string) func(*lua.LState) int {
	return func(L *lua.LState) int {
		original := L.ToString(1)
		key := util.GetTranslationKey(original)
		translation, ok := dictionary[key]
		if !ok {
			translation = original
		}
		L.Push(lua.LString(translation))
		return 1
	}
}

func jsonDecode(L *lua.LState) int {
	jsonStr := L.ToString(1)
	var result interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(convertToLuaValue(L, result))
	return 1
}

func jsonEncode(L *lua.LState) int {
	luaValue := L.Get(1)
	jsonBytes, err := json.Marshal(convertFromLuaValue(luaValue))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(jsonBytes)))
	return 1
}

func convertToLuaValue(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		table := L.NewTable()
		for i, item := range v {
			table.RawSetInt(i+1, convertToLuaValue(L, item))
		}
		return table
	case map[string]interface{}:
		table := L.NewTable()
		for k, val := range v {
			table.RawSetString(k, convertToLuaValue(L, val))
		}
		return table
	default:
		return lua.LNil
	}
}

func convertFromLuaValue(lv lua.LValue) interface{} {
	switch lv.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return bool(lv.(lua.LBool))
	case lua.LTNumber:
		return float64(lv.(lua.LNumber))
	case lua.LTString:
		return string(lv.(lua.LString))
	case lua.LTTable:
		table := lv.(*lua.LTable)
		result := make(map[string]interface{})
		isArray := true
		maxIndex := 0

		table.ForEach(func(key, value lua.LValue) {
			if num, ok := key.(lua.LNumber); ok {
				idx := int(num)
				if idx > maxIndex {
					maxIndex = idx
				}
				if idx < 1 || idx > maxIndex {
					isArray = false
				}
			} else {
				isArray = false
			}
		})

		if isArray && maxIndex > 0 {
			arr := make([]interface{}, maxIndex)
			table.ForEach(func(key, value lua.LValue) {
				if num, ok := key.(lua.LNumber); ok {
					idx := int(num) - 1
					if idx >= 0 && idx < maxIndex {
						arr[idx] = convertFromLuaValue(value)
					}
				}
			})
			return arr
		}

		table.ForEach(func(key, value lua.LValue) {
			keyStr := key.String()
			result[keyStr] = convertFromLuaValue(value)
		})
		return result
	default:
		return nil
	}
}
