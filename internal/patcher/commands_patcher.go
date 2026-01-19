package patcher

import (
	"encoding/json"
	"htpatcher/internal/domain"
	"htpatcher/internal/domain/rpgmaker"
	"htpatcher/internal/util"
	"slices"
	"strings"
)

// patchVariableValue handles translation of variable assignment values
// Supports double-quoted strings, single-quoted strings, and single-quoted JSON arrays
func patchVariableValue(value string, dictionary map[string]string) string {
	// Handle double-quoted strings: "text"
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		s := value[1 : len(value)-1]
		if translation, ok := dictionary[util.GetTranslationKey(s)]; ok {
			return "\"" + translation + "\""
		}
		return value
	}

	// Handle double-quoted strings with semicolon: "text";
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\";") {
		s := value[1 : len(value)-2]
		if translation, ok := dictionary[util.GetTranslationKey(s)]; ok {
			return "\"" + translation + "\";"
		}
		return value
	}

	// Handle single-quoted strings: 'text' or '[...]'
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		s := value[1 : len(value)-1]

		// Check if it's a JSON array: '[...]'
		if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
			var jsonArray []any
			if err := json.Unmarshal([]byte(s), &jsonArray); err == nil {
				newArray := make([]any, 0, len(jsonArray))
				for _, item := range jsonArray {
					if str, ok := item.(string); ok {
						if translation, ok := dictionary[util.GetTranslationKey(str)]; ok {
							newArray = append(newArray, translation)
						} else {
							newArray = append(newArray, str)
						}
					} else {
						newArray = append(newArray, item)
					}
				}
				if jsonBytes, err := json.Marshal(newArray); err == nil {
					return "'" + string(jsonBytes) + "'"
				}
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

// patchParameterValue recursively traverses data structures and applies translations to strings
func patchParameterValue(value any, dictionary map[string]string) any {
	switch v := value.(type) {
	case string:
		// Base case: translate string if found in dictionary
		if translation, ok := dictionary[util.GetTranslationKey(v)]; ok {
			return translation
		}
		return v
	case []any:
		// Recursive case: process array elements
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = patchParameterValue(item, dictionary)
		}
		return result
	case map[string]any:
		// Recursive case: process object properties
		result := make(map[string]any)
		for key, val := range v {
			result[key] = patchParameterValue(val, dictionary)
		}
		return result
	case *util.OrderedMap:
		// Recursive case: process ordered map properties while preserving key order
		result := util.NewOrderedMap()
		for _, key := range v.Keys {
			result.Set(key, patchParameterValue(v.Values[key], dictionary))
		}
		return result
	default:
		// Other types (numbers, booleans, null): return as-is
		return v
	}
}

// patchCommands patches event commands
func patchCommands(commands []*rpgmaker.EventCommand, patchInfo *domain.PatchInfo) ([]*rpgmaker.EventCommand, error) {
	commandsToDelete := []int{}
	commandIndex := 0
	last101CommandHasSpeakerThumbnail := false

	for commandIndex < len(commands) {
		command := commands[commandIndex]

		// Command 101 is start of dialogue, and if param 4 is a string, it is the speaker name
		if command.Code == 101 {
			if len(command.Parameters) > 4 {
				if key, ok := command.Parameters[4].(string); ok {
					if speakerName, ok := patchInfo.Dictionary[util.GetTranslationKey(key)]; ok {
						command.Parameters[4] = speakerName
					}
				}
			}
			last101CommandHasSpeakerThumbnail = false
			if thumbnail, ok := command.Parameters[0].(string); ok {
				if thumbnail != "" {
					last101CommandHasSpeakerThumbnail = true
				}
			}
		}

		// Command 401 is continuation of dialogue and param 0 is the text
		if command.Code == 401 {
			wrapWidth := patchInfo.Config.WrapWidth
			if last101CommandHasSpeakerThumbnail && patchInfo.Config.DynamicWrapWidth {
				wrapWidth -= 10
			}

			dialogueCommands := []*rpgmaker.EventCommand{}
			fullText := ""

			// Collect all consecutive 401 commands
			for commandIndex < len(commands) && commands[commandIndex].Code == 401 {
				dialogueCommands = append(dialogueCommands, commands[commandIndex])
				if text, ok := commands[commandIndex].Parameters[0].(string); ok {
					fullText += text
				}
				commandIndex++
			}
			commandIndex--

			translationKey := util.GetTranslationKey(fullText)
			if translation, ok := patchInfo.Dictionary[translationKey]; ok {
				dialogueCommands[0].Parameters[0] = util.Wrap(translation, wrapWidth)
				// Only keep the first command in the dialogue
				for k := (commandIndex - len(dialogueCommands) + 2); k <= commandIndex; k++ {
					commandsToDelete = append(commandsToDelete, k)
				}
			}
		}

		// Command 405 is rolling text and param 0 is the text
		if command.Code == 405 {
			if text, ok := command.Parameters[0].(string); ok {
				if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(text)]; ok {
					command.Parameters[0] = translation
				}
			}
		}

		// Command 102 is display choices and param 0 is an array of string choices
		if command.Code == 102 {
			if choices, ok := command.Parameters[0].([]any); ok {
				for i, choice := range choices {
					if choice, ok := choice.(string); ok {
						if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(choice)]; ok {
							command.Parameters[0].([]any)[i] = translation
						}
					}
				}
			}
		}

		// Command 408 is choice description and param 0 is the description
		if command.Code == 408 {
			if description, ok := command.Parameters[0].(string); ok {
				if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(description)]; ok {
					command.Parameters[0] = util.Wrap(util.NoNewline(translation), patchInfo.Config.WrapWidth)
				}
			}
		}

		// Command 303 is name input and param 1 is the length of the input
		if command.Code == 303 {
			length, ok := command.Parameters[1].(float64)
			if ok && length < 10 {
				command.Parameters[1] = 10
			}
		}

		// Command 122 is variable assignment and param 4 is the value
		if command.Code == 122 {
			if len(command.Parameters) > 4 {
				// Check if param 0 is a variable ID that needs patching
				if varID, ok := command.Parameters[0].(float64); ok {
					if slices.Contains(patchInfo.Config.VariablesToPatch, int(varID)) {
						// Check if param 4 is a string value
						if value, ok := command.Parameters[4].(string); ok {
							command.Parameters[4] = patchVariableValue(value, patchInfo.Dictionary)
						}
					}
				}
			}
		}

		// Command 355 is Script, 655 is Script continuation
		// Join all lines with \n, translate, put result in 355, delete 655s
		if command.Code == 355 {
			fullScript := ""

			if text, ok := command.Parameters[0].(string); ok {
				fullScript = text
			}

			// Collect consecutive 655 commands, joining with \n
			startOfContinuation := commandIndex + 1
			nextIndex := startOfContinuation
			for nextIndex < len(commands) && commands[nextIndex].Code == 655 {
				if text, ok := commands[nextIndex].Parameters[0].(string); ok {
					fullScript += "\n" + text
				}
				nextIndex++
			}

			// Look up translation
			if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(fullScript)]; ok {
				// Put entire translation in 355 command
				command.Parameters[0] = translation

				// Mark all 655 commands for deletion
				for i := startOfContinuation; i < nextIndex; i++ {
					commandsToDelete = append(commandsToDelete, i)
				}
			}

			commandIndex = nextIndex - 1
		}

		// Command 357 is a plugin call, param 0 is the plugin name, param 1 is the function name, param 3 is the options
		if command.Code == 357 {
			if len(command.Parameters) > 3 {
				if plugin, ok := command.Parameters[0].(string); ok {
					if function, ok := command.Parameters[1].(string); ok {
						for _, parameter := range patchInfo.Config.ParametersToPatch {
							if parameter.Plugin == plugin && parameter.Function == function {
								switch parameter.RootType {
								case "string":
									if options, ok := command.Parameters[3].(string); ok {
										if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(options)]; ok {
											command.Parameters[3] = util.Wrap(translation, patchInfo.Config.WrapWidth)
										}
									}
								case "array":
									if options, ok := command.Parameters[3].([]any); ok {
										command.Parameters[3] = patchParameterValue(options, patchInfo.Dictionary)
									}
								case "object":
									command.Parameters[3] = patchParameterValue(command.Parameters[3], patchInfo.Dictionary)
								}
							}
						}
					}
				}
			}
		}

		commandIndex++
	}

	// Filter out commands marked for deletion
	newCommands := commands[:0]
	for idx, command := range commands {
		if slices.Contains(commandsToDelete, idx) {
			continue
		}
		newCommands = append(newCommands, command)
	}

	return newCommands, nil
}
