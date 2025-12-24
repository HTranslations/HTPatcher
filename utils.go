package main

import (
	"encoding/json"
	"strings"
)

func Wrap(text string, width int) string {
	if width <= 0 {
		width = 58
	}

	if text == "" || len(text) <= width {
		return text
	}

	words := strings.Split(text, " ")
	lines := []string{}
	currentLine := ""

	for _, word := range words {
		// If adding this word would exceed the width, start a new line
		if len(currentLine)+len(word)+1 > width {
			if len(currentLine) > 0 {
				lines = append(lines, strings.TrimSpace(currentLine))
				currentLine = word
			} else {
				// If the word itself is longer than width, break it
				if len(word) > width {
					// Break long words at character boundaries
					for i := 0; i < len(word); i += width {
						end := min(i+width, len(word))
						lines = append(lines, word[i:end])
					}
				} else {
					lines = append(lines, word)
				}
			}
		} else {
			// Add word to current line
			if len(currentLine) > 0 {
				currentLine += " " + word
			} else {
				currentLine = word
			}
		}
	}

	// Add the last line if it has content
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
