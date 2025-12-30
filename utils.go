package main

import (
	"encoding/json"
	"regexp"
	"strings"
)

var placeholderRegex = regexp.MustCompile(`\\[A-Za-z]*(\[[^\]]*\])?`)

func visibleLength(text string) int {
	matches := placeholderRegex.FindAllString(text, -1)

	totalPlaceholderLength := 0
	for _, match := range matches {
		if !strings.HasPrefix(match, "\\N[") {
			totalPlaceholderLength += len(match)
		}
	}

	return len(text) - totalPlaceholderLength
}

func Wrap(text string, width int) string {
	if width <= 0 {
		width = 58
	}

	if text == "" || visibleLength(text) <= width {
		return text
	}

	words := strings.Split(text, " ")
	lines := []string{}
	currentLine := ""

	for _, word := range words {
		spaceLength := 0
		if len(currentLine) > 0 {
			spaceLength = 1
		}
		currentVisibleLen := visibleLength(currentLine)
		wordVisibleLen := visibleLength(word)

		if currentVisibleLen+wordVisibleLen+spaceLength > width {
			if len(currentLine) > 0 {
				lines = append(lines, strings.TrimSpace(currentLine))
				currentLine = word
			} else {
				if wordVisibleLen > width {
					visibleChars := 0
					actualPos := 0
					wordStart := 0

					for actualPos < len(word) {
						loc := placeholderRegex.FindStringIndex(word[actualPos:])

						if loc != nil && loc[0] == 0 {
							placeholder := word[actualPos : actualPos+loc[1]]
							actualPos += loc[1]

							if strings.HasPrefix(placeholder, "\\N[") {
								visibleChars += len(placeholder)
							}
						} else {
							visibleChars++
							actualPos++
						}

						if visibleChars >= width {
							lines = append(lines, word[wordStart:actualPos])
							wordStart = actualPos
							visibleChars = 0
						}
					}

					if wordStart < len(word) {
						lines = append(lines, word[wordStart:])
					}
				} else {
					lines = append(lines, word)
				}
			}
		} else {
			if len(currentLine) > 0 {
				currentLine += " " + word
			} else {
				currentLine = word
			}
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, strings.TrimSpace(currentLine))
	}

	return strings.Join(lines, "\n")
}

func NoNewline(text string) string {
	return strings.ReplaceAll(text, "\n", " ")
}

// MergeJsonChanges merges typed changes into raw JSON while preserving unknown fields.
// Supports both top-level maps and arrays.
func MergeJsonChanges(data []byte, changes any) ([]byte, error) {
	var original any
	if err := json.Unmarshal(data, &original); err != nil {
		return nil, err
	}

	changedBytes, err := json.Marshal(changes) // no & needed
	if err != nil {
		return nil, err
	}

	var changed any
	if err := json.Unmarshal(changedBytes, &changed); err != nil {
		return nil, err
	}

	merged := mergeAny(original, changed, "")
	return json.Marshal(merged)
}

func mergeAny(original, changes any, key string) any {
	// Special rule: "list" is always replaced, never merged
	if key == "list" {
		return changes
	}

	switch o := original.(type) {

	case map[string]any:
		if c, ok := changes.(map[string]any); ok {
			for k, v := range c {
				if ov, exists := o[k]; exists {
					o[k] = mergeAny(ov, v, k)
				} else {
					o[k] = v
				}
			}
			return o
		}
		return changes

	case []any:
		if c, ok := changes.([]any); ok {
			minLen := len(o)
			if len(c) < minLen {
				minLen = len(c)
			}
			for i := 0; i < minLen; i++ {
				o[i] = mergeAny(o[i], c[i], "")
			}
			if len(c) > len(o) {
				o = append(o, c[len(o):]...)
			}
			return o
		}
		return changes

	default:
		// primitives
		return changes
	}
}

func GetTranslationKey(text string) string {
	return strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), " ", ""))
}
