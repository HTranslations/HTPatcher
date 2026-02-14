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

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", patchInfo.Dictionary, false, 0)
		patchStringProperty(obj, "nickname", patchInfo.Dictionary, false, 0)
		patchStringProperty(obj, "description", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", patchInfo.Dictionary, false, 0)
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

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", patchInfo.Dictionary, false, 0)
		patchStringProperty(obj, "note", patchInfo.Dictionary, false, 0)
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

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", patchInfo.Dictionary, false, 0)
		patchStringProperty(obj, "description", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "message1", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "message2", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", patchInfo.Dictionary, false, 0)
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

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", patchInfo.Dictionary, false, 0)
		patchStringProperty(obj, "description", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", patchInfo.Dictionary, false, 0)
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

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", patchInfo.Dictionary, false, 0)
		patchStringProperty(obj, "description", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", patchInfo.Dictionary, false, 0)
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

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", patchInfo.Dictionary, false, 0)
		patchStringProperty(obj, "description", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", patchInfo.Dictionary, false, 0)
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

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", patchInfo.Dictionary, false, 0)
		patchStringProperty(obj, "note", patchInfo.Dictionary, false, 0)
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

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", patchInfo.Dictionary, false, 0)
		patchStringProperty(obj, "message1", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "message2", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "message3", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "message4", patchInfo.Dictionary, true, patchInfo.Config.WrapWidth)
		patchStringProperty(obj, "note", patchInfo.Dictionary, false, 0)
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

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		patchStringProperty(obj, "name", patchInfo.Dictionary, false, 0)

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

	for _, item := range arr {
		if item == nil {
			continue
		}
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

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

	// Patch display name
	patchStringProperty(obj, "display_name", patchInfo.Dictionary, false, 0)

	// Patch events - VX Ace uses a hash/map, not an array
	eventsMap := obj.GetMap("events")
	if eventsMap != nil {
		for _, eventItem := range eventsMap {
			eventObj, ok := eventItem.(*marshal.RubyObject)
			if !ok {
				continue
			}

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

	// Patch game title
	patchStringProperty(obj, "game_title", patchInfo.Dictionary, false, 0)

	// Patch currency unit
	patchStringProperty(obj, "currency_unit", patchInfo.Dictionary, false, 0)

	// Patch string arrays
	patchStringArray(obj, "elements", patchInfo.Dictionary)
	patchStringArray(obj, "skill_types", patchInfo.Dictionary)
	patchStringArray(obj, "weapon_types", patchInfo.Dictionary)
	patchStringArray(obj, "armor_types", patchInfo.Dictionary)
	patchStringArray(obj, "equip_types", patchInfo.Dictionary)
	patchStringArray(obj, "switches", patchInfo.Dictionary)
	patchStringArray(obj, "variables", patchInfo.Dictionary)

	// Patch terms object
	termsObj := obj.GetObject("terms")
	if termsObj != nil {
		patchStringArray(termsObj, "basic", patchInfo.Dictionary)
		patchStringArray(termsObj, "params", patchInfo.Dictionary)
		patchStringArray(termsObj, "etypes", patchInfo.Dictionary)
		patchStringArray(termsObj, "commands", patchInfo.Dictionary)
	}

	return marshal.Write(obj)
}

// Helper functions

// patchStringProperty patches a string property on a RubyObject
func patchStringProperty(obj *marshal.RubyObject, key string, dictionary map[string]string, wrap bool, wrapWidth int) {
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

	translationKey := util.GetTranslationKey(str)
	if translation, exists := dictionary[translationKey]; exists {
		if wrap && wrapWidth > 0 {
			translation = util.Wrap(util.NoNewline(translation), wrapWidth)
		}
		obj.Properties[key] = translation
	}
}

// patchStringArray patches translatable strings in an array property
func patchStringArray(obj *marshal.RubyObject, key string, dictionary map[string]string) {
	if obj == nil || obj.Properties == nil {
		return
	}

	arr := obj.GetArray(key)
	if arr == nil {
		return
	}

	for i, v := range arr {
		if s, ok := v.(string); ok {
			translationKey := util.GetTranslationKey(s)
			if translation, exists := dictionary[translationKey]; exists {
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

	for commandIndex < len(commands) {
		command := commands[commandIndex]
		code := command.getCode()
		params := command.getParameters()

		// Command 101: Start of dialogue
		if code == 101 {
			if len(params) > 4 {
				if key, ok := params[4].(string); ok {
					translationKey := util.GetTranslationKey(key)
					if translation, exists := patchInfo.Dictionary[translationKey]; exists {
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

			translationKey := util.GetTranslationKey(fullText)
			if translation, exists := patchInfo.Dictionary[translationKey]; exists {
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

			translationKey := util.GetTranslationKey(fullText)
			if translation, exists := patchInfo.Dictionary[translationKey]; exists {
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
							if translation, exists := patchInfo.Dictionary[util.GetTranslationKey(choiceStr)]; exists {
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
					if translation, exists := patchInfo.Dictionary[util.GetTranslationKey(text)]; exists {
						command.setParameter(0, translation)
					}
				}
			}
		}

		// Command 408: Choice description
		if code == 408 {
			if len(params) > 0 {
				if description, ok := params[0].(string); ok {
					if translation, exists := patchInfo.Dictionary[util.GetTranslationKey(description)]; exists {
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
			if translation, exists := patchInfo.Dictionary[util.GetTranslationKey(fullScript)]; exists {
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
								patched := patchVariableValue(value, patchInfo.Dictionary)
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
func patchVariableValue(value string, dictionary map[string]string) string {
	// Handle double-quoted strings: "text"
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		s := value[1 : len(value)-1]
		if translation, ok := dictionary[util.GetTranslationKey(s)]; ok {
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
			if translation, ok := dictionary[util.GetTranslationKey(s)]; ok {
				return "'" + translation + "'"
			}
			return value
		}
		// Regular single-quoted string
		if translation, ok := dictionary[util.GetTranslationKey(s)]; ok {
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
		translationKey := util.GetTranslationKey(sourceStr)
		if translation, exists := patchInfo.Dictionary[translationKey]; exists {
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
