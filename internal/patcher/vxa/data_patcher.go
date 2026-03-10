package vxa

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/marshal"
	"htpatcher/internal/util"
	"io"
	"strconv"
)

// patchActors patches Actors.rvdata2
func patchActors(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", dict, km, "actor_name", false, 0)
		patchStringProperty(obj, "nickname", dict, km, "actor_nickname", false, 0)
		patchStringProperty(obj, "description", dict, km, "actor_profile", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", dict, km, "actor_note", false, 0)
	}

	return marshal.Write(arr)
}

// patchClasses patches Classes.rvdata2
func patchClasses(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", dict, km, "class_name", false, 0)
		patchStringProperty(obj, "note", dict, km, "class_name", false, 0)
	}

	return marshal.Write(arr)
}

// patchSkills patches Skills.rvdata2
func patchSkills(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", dict, km, "skill_name", false, 0)
		patchStringProperty(obj, "description", dict, km, "skill_description", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "message1", dict, km, "skill_message1", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "message2", dict, km, "skill_message2", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", dict, km, "skill_name", false, 0)
	}

	return marshal.Write(arr)
}

// patchItems patches Items.rvdata2
func patchItems(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", dict, km, "item_name", false, 0)
		patchStringProperty(obj, "description", dict, km, "item_description", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", dict, km, "item_note", false, 0)
	}

	return marshal.Write(arr)
}

// patchWeapons patches Weapons.rvdata2
func patchWeapons(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", dict, km, "weapon_name", false, 0)
		patchStringProperty(obj, "description", dict, km, "weapon_description", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", dict, km, "weapon_name", false, 0)
	}

	return marshal.Write(arr)
}

// patchArmors patches Armors.rvdata2
func patchArmors(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", dict, km, "armor_name", false, 0)
		patchStringProperty(obj, "description", dict, km, "armor_description", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", dict, km, "armor_name", false, 0)
	}

	return marshal.Write(arr)
}

// patchEnemies patches Enemies.rvdata2
func patchEnemies(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", dict, km, "enemy_name", false, 0)
		patchStringProperty(obj, "note", dict, km, "enemy_note", false, 0)
	}

	return marshal.Write(arr)
}

// patchStates patches States.rvdata2
func patchStates(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", dict, km, "state_name", false, 0)
		patchStringProperty(obj, "message1", dict, km, "state_message1", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "message2", dict, km, "state_message2", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "message3", dict, km, "state_message3", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "message4", dict, km, "state_message4", true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", dict, km, "state_note", false, 0)
	}

	return marshal.Write(arr)
}

// patchTroops patches Troops.rvdata2
func patchTroops(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", dict, km, "troop_name", false, 0)

		// Patch pages
		pages := obj.GetArray("pages")
		for _, pageItem := range pages {
			pageObj, ok := pageItem.(*marshal.RubyObject)
			if !ok {
				continue
			}
			patchEventCommands(pageObj, patchInfo)
		}
	}

	return marshal.Write(arr)
}

// patchCommonEvents patches CommonEvents.rvdata2
func patchCommonEvents(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", dict, km, "common_event_name", false, 0)
		patchEventCommands(obj, patchInfo)
	}

	return marshal.Write(arr)
}

// patchMap patches a Map*.rvdata2 file
func patchMap(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	obj, ok := raw.(*marshal.RubyObject)
	if !ok {
		return nil, fmt.Errorf("expected RubyObject, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	// Patch display name
	patchStringProperty(obj, "display_name", dict, km, "map_display_name", false, 0)

	// Patch events - VX Ace uses a hash/map, not an array
	eventsMap := obj.GetMap("events")
	if eventsMap != nil {
		for _, eventItem := range eventsMap {
			eventObj, ok := eventItem.(*marshal.RubyObject)
			if !ok {
				continue
			}

			// Patch event name
			patchStringProperty(eventObj, "name", dict, km, "event_name", false, 0)

			// Patch event pages
			pages := eventObj.GetArray("pages")
			for _, pageItem := range pages {
				pageObj, ok := pageItem.(*marshal.RubyObject)
				if !ok {
					continue
				}
				patchEventCommands(pageObj, patchInfo)
			}
		}
	}

	return marshal.Write(obj)
}

// patchSystem patches System.rvdata2
func patchSystem(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	obj, ok := raw.(*marshal.RubyObject)
	if !ok {
		return nil, fmt.Errorf("expected RubyObject, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	// Patch game title
	patchStringProperty(obj, "game_title", dict, km, "game_title", false, 0)

	// Patch currency unit
	patchStringProperty(obj, "currency_unit", dict, km, "system_term", false, 0)

	// Patch string arrays
	patchStringArray(obj, "elements", dict, km, "system_term")
	patchStringArray(obj, "skill_types", dict, km, "system_term")
	patchStringArray(obj, "weapon_types", dict, km, "system_term")
	patchStringArray(obj, "armor_types", dict, km, "system_term")
	patchStringArray(obj, "equip_types", dict, km, "system_term")
	patchStringArray(obj, "switches", dict, km, "system_switch")
	patchStringArray(obj, "variables", dict, km, "system_variable")

	// Patch terms object
	termsObj := obj.GetObject("terms")
	if termsObj != nil {
		patchStringArray(termsObj, "basic", dict, km, "system_term")
		patchStringArray(termsObj, "params", dict, km, "system_term")
		patchStringArray(termsObj, "etypes", dict, km, "system_term")
		patchStringArray(termsObj, "commands", dict, km, "system_term")
	}

	return marshal.Write(obj)
}

// Helper functions

// patchStringProperty patches a string property on a RubyObject
func patchStringProperty(obj *marshal.RubyObject, key string, dictionary map[string]string, keyMode string, entryType string, wrap bool, wrapWidth int) {
	if obj == nil || obj.Properties == nil {
		return
	}

	val, ok := obj.Properties[key]
	if !ok {
		return
	}

	str, ok := val.(string)
	if !ok {
		return
	}

	if translation, exists := util.DictLookup(dictionary, keyMode, entryType, str); exists {
		if wrap && wrapWidth > 0 {
			translation = util.Wrap(util.NoNewline(translation), wrapWidth)
		}
		obj.Properties[key] = translation
	}
}

// patchStringArray patches translatable strings in an array property
func patchStringArray(obj *marshal.RubyObject, key string, dictionary map[string]string, keyMode string, entryType string) {
	if obj == nil || obj.Properties == nil {
		return
	}

	arr := obj.GetArray(key)
	if arr == nil {
		return
	}

	for i, v := range arr {
		if s, ok := v.(string); ok {
			if translation, exists := util.DictLookup(dictionary, keyMode, entryType, s); exists {
				arr[i] = translation
			}
		}
	}
}

// patchEventCommands patches the event command list on a RubyObject that has a "list" property
func patchEventCommands(obj *marshal.RubyObject, patchInfo *domain.PatchInfo) {
	if obj == nil || obj.Properties == nil {
		return
	}

	listRaw := obj.GetArray("list")
	if listRaw == nil {
		return
	}

	// Convert to command objects for patching
	commands := make([]*commandWrapper, 0, len(listRaw))
	for _, item := range listRaw {
		cmdObj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}
		commands = append(commands, &commandWrapper{obj: cmdObj})
	}

	// Patch commands
	newCommands := patchCommands(commands, patchInfo)

	// Convert back to interface array
	newList := make([]interface{}, len(newCommands))
	for i, cmd := range newCommands {
		newList[i] = cmd.obj
	}

	obj.Properties["list"] = newList
}

// commandWrapper wraps a RubyObject for command patching
type commandWrapper struct {
	obj *marshal.RubyObject
}

func (c *commandWrapper) getCode() int {
	return c.obj.GetInt("code")
}

func (c *commandWrapper) getIndent() int {
	return c.obj.GetInt("indent")
}

func (c *commandWrapper) getParameters() []interface{} {
	return c.obj.GetArray("parameters")
}

func (c *commandWrapper) setParameter(idx int, value interface{}) {
	params := c.obj.GetArray("parameters")
	if params != nil && idx < len(params) {
		params[idx] = value
	}
}

func (c *commandWrapper) setParameters(params []interface{}) {
	c.obj.Properties["parameters"] = params
}

// patchCommands patches a list of event commands
func patchCommands(commands []*commandWrapper, patchInfo *domain.PatchInfo) []*commandWrapper {
	commandsToDelete := make(map[int]bool)
	commandIndex := 0
	last101HasThumbnail := false

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for commandIndex < len(commands) {
		command := commands[commandIndex]
		code := command.getCode()
		params := command.getParameters()

		// Command 101: Start of dialogue
		if code == 101 {
			if len(params) > 4 {
				if key, ok := params[4].(string); ok {
					if translation, exists := util.DictLookup(dict, km, "speaker", key); exists {
						command.setParameter(4, translation)
					}
				}
			}
			last101HasThumbnail = false
			if len(params) > 0 {
				if thumbnail, ok := params[0].(string); ok && thumbnail != "" {
					last101HasThumbnail = true
				}
			}
		}

		// Command 401: Dialogue text continuation
		if code == 401 {
			wrapWidth := patchInfo.Config.WrapWidth
			if last101HasThumbnail && patchInfo.Config.DynamicWrapWidth {
				wrapWidth -= 15
			}

			// Collect all consecutive 401 commands
			dialogueCommands := []*commandWrapper{}
			fullText := ""

			for commandIndex < len(commands) && commands[commandIndex].getCode() == 401 {
				dialogueCommands = append(dialogueCommands, commands[commandIndex])
				dialogueParams := commands[commandIndex].getParameters()
				if len(dialogueParams) > 0 {
					if text, ok := dialogueParams[0].(string); ok {
						fullText += text
					}
				}
				commandIndex++
			}
			commandIndex--

			if translation, exists := util.DictLookup(dict, km, "dialogue", fullText); exists {
				dialogueCommands[0].setParameter(0, util.Wrap(translation, wrapWidth))
				// Mark subsequent 401 commands for deletion
				startIdx := commandIndex - len(dialogueCommands) + 2
				for k := startIdx; k <= commandIndex; k++ {
					commandsToDelete[k] = true
				}
			}
		}

		// Command 405: Scrolling text
		if code == 405 {
			scrollingCommands := []*commandWrapper{}
			fullText := ""

			// Collect all consecutive 405 commands
			for commandIndex < len(commands) && commands[commandIndex].getCode() == 405 {
				scrollingCommands = append(scrollingCommands, commands[commandIndex])
				scrollingParams := commands[commandIndex].getParameters()
				if len(scrollingParams) > 0 {
					if text, ok := scrollingParams[0].(string); ok {
						fullText += text
					}
				}
				commandIndex++
			}
			commandIndex--

			if translation, exists := util.DictLookup(dict, km, "message", fullText); exists {
				scrollingCommands[0].setParameter(0, util.Wrap(translation, patchInfo.Config.WrapWidth))
				// Mark subsequent 405 commands for deletion
				startIdx := commandIndex - len(scrollingCommands) + 2
				for k := startIdx; k <= commandIndex; k++ {
					commandsToDelete[k] = true
				}
			}
		}

		// Command 102: Show choices
		if code == 102 {
			if len(params) > 0 {
				if choices, ok := params[0].([]interface{}); ok {
					for i, choice := range choices {
						if choiceStr, ok := choice.(string); ok {
							if translation, exists := util.DictLookup(dict, km, "choice", choiceStr); exists {
								choices[i] = translation
							}
						}
					}
				}
			}
		}

		// Command 108: Comment
		if code == 108 {
			if len(params) > 0 {
				if text, ok := params[0].(string); ok {
					if translation, exists := util.DictLookup(dict, km, "comment", text); exists {
						command.setParameter(0, translation)
					}
				}
			}
		}

		// Command 408: Choice description
		if code == 408 {
			if len(params) > 0 {
				if description, ok := params[0].(string); ok {
					if translation, exists := util.DictLookup(dict, km, "comment", description); exists {
						command.setParameter(0, util.Wrap(util.NoNewline(translation), patchInfo.Config.WrapWidth))
					}
				}
			}
		}

		// Command 303: Name input processing
		if code == 303 {
			if len(params) > 1 {
				if length, ok := params[1].(int); ok && length < 10 {
					command.setParameter(1, 10)
				}
			}
		}

		// Command 355: Script, 655: Script continuation
		if code == 355 {
			fullScript := ""
			if len(params) > 0 {
				if text, ok := params[0].(string); ok {
					fullScript = text
				}
			}

			// Collect consecutive 655 commands
			startOfContinuation := commandIndex + 1
			nextIndex := startOfContinuation
			for nextIndex < len(commands) && commands[nextIndex].getCode() == 655 {
				nextParams := commands[nextIndex].getParameters()
				if len(nextParams) > 0 {
					if text, ok := nextParams[0].(string); ok {
						fullScript += "\n" + text
					}
				}
				nextIndex++
			}

			// Look up translation
			if translation, exists := util.DictLookup(dict, km, "script", fullScript); exists {
				command.setParameter(0, translation)
				// Mark all 655 commands for deletion
				for i := startOfContinuation; i < nextIndex; i++ {
					commandsToDelete[i] = true
				}
			}

			commandIndex = nextIndex - 1
		}

		// Command 122: Variable assignment (for plugin-specific translations)
		if code == 122 && len(patchInfo.Config.VariablesToPatch) > 0 {
			if len(params) > 4 {
				if varID, ok := params[0].(int); ok {
					for _, patchVarID := range patchInfo.Config.VariablesToPatch {
						if varID == patchVarID {
							if value, ok := params[4].(string); ok {
								patched := patchVariableValue(value, dict, km)
								command.setParameter(4, patched)
							}
							break
						}
					}
				}
			}
		}

		commandIndex++
	}

	// Filter out deleted commands
	result := make([]*commandWrapper, 0, len(commands))
	for i, cmd := range commands {
		if !commandsToDelete[i] {
			result = append(result, cmd)
		}
	}

	return result
}

// patchVariableValue patches variable assignment values (for plugin-specific text)
func patchVariableValue(value string, dictionary map[string]string, keyMode string) string {
	// Handle double-quoted strings: "text"
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		s := value[1 : len(value)-1]
		if translation, ok := util.DictLookup(dictionary, keyMode, "variable_value", s); ok {
			return "\"" + translation + "\""
		}
		return value
	}

	// Handle single-quoted strings: 'text'
	if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
		s := value[1 : len(value)-1]
		// Check if it's a JSON array
		if len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']' {
			// Parse JSON array and translate each string element
			// Simplified handling - just translate the whole thing if found
			if translation, ok := util.DictLookup(dictionary, keyMode, "variable_value", s); ok {
				return "'" + translation + "'"
			}
			return value
		}
		// Regular single-quoted string
		if translation, ok := util.DictLookup(dictionary, keyMode, "variable_value", s); ok {
			return "'" + translation + "'"
		}
		return value
	}

	return value
}

// getIntFromInterface converts an interface to int
func getIntFromInterface(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case string:
		if i, err := strconv.Atoi(n); err == nil {
			return i
		}
	}
	return 0
}

// patchScripts patches Scripts.rvdata2
// Scripts are stored as an array of [section_id, name, zlib_compressed_source]
func patchScripts(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	raw, err := marshal.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range arr {
		entry, ok := item.([]interface{})
		if !ok || len(entry) < 3 {
			continue
		}

		// Get compressed source (3rd element)
		compressedStr, ok := entry[2].(string)
		if !ok {
			continue
		}

		// Decompress
		source, err := decompressZlib([]byte(compressedStr))
		if err != nil {
			continue // Skip scripts that can't be decompressed
		}

		sourceStr := string(source)
		if sourceStr == "" {
			continue
		}

		// Look up in dictionary
		if translation, exists := util.DictLookup(dict, km, "script", sourceStr); exists {
			// Recompress the translated source
			compressed, err := compressZlib([]byte(translation))
			if err != nil {
				return nil, fmt.Errorf("failed to compress translated script: %w", err)
			}
			entry[2] = string(compressed)
		}
	}

	return marshal.Write(arr)
}

// decompressZlib decompresses zlib-compressed data
func decompressZlib(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

// compressZlib compresses data using zlib
func compressZlib(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
