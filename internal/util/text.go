package util

import (
	"regexp"
	"strings"
	"unicode"
)

var placeholderRegex = regexp.MustCompile(`\\[A-Za-z]*(\[[^\]]*\])?`)

// VisibleLength calculates the visible length of text excluding placeholders
func VisibleLength(text string) int {
	matches := placeholderRegex.FindAllString(text, -1)

	totalPlaceholderLength := 0
	for _, match := range matches {
		if !strings.HasPrefix(match, "\\N[") {
			totalPlaceholderLength += len(match)
		}
	}

	return len(text) - totalPlaceholderLength
}

// Wrap wraps text to a specified width, accounting for RPG Maker placeholders.
// Existing newlines are preserved and reset the line width counter.
func Wrap(text string, width int) string {
	if width <= 0 {
		width = 58
	}

	if text == "" {
		return text
	}

	segments := strings.Split(text, "\n")
	wrappedSegments := make([]string, 0, len(segments))

	for _, segment := range segments {
		wrappedSegments = append(wrappedSegments, wrapSegment(segment, width))
	}

	return strings.Join(wrappedSegments, "\n")
}

// wrapSegment wraps a single line of text (no newlines) to the specified width.
func wrapSegment(text string, width int) string {
	if text == "" || VisibleLength(text) <= width {
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
		currentVisibleLen := VisibleLength(currentLine)
		wordVisibleLen := VisibleLength(word)

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

// NoNewline removes newlines from text, replacing them with spaces
func NoNewline(text string) string {
	return strings.ReplaceAll(text, "\n", " ")
}

// GetTranslationKey generates a normalized key for dictionary lookup (legacy, untyped)
func GetTranslationKey(text string) string {
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\u3000", "") // full-width space
	return strings.ToLower(text)
}

// GetTypedTranslationKey generates a type-prefixed key for dictionary lookup
func GetTypedTranslationKey(entryType string, text string) string {
	normalized := GetTranslationKey(text)
	if normalized == "" {
		return ""
	}
	return entryType + ":" + normalized
}

// IsJapaneseChar checks if a single rune is a Japanese character
func IsJapaneseChar(r rune) bool {
	return unicode.In(r,
		unicode.Hiragana,
		unicode.Katakana,
		unicode.Han,
	) ||
		(r >= 0x30A0 && r <= 0x30FF) ||
		(r >= 0xFF65 && r <= 0xFF9F) ||
		(r >= 0x3000 && r <= 0x303F)
}

// ContainsJapanese checks if a string contains Japanese characters
func ContainsJapanese(s string) bool {
	for _, r := range s {
		if IsJapaneseChar(r) {
			return true
		}
	}
	return false
}

// IsOnlyPunctuation checks if a string contains only punctuation, whitespace and symbols
func IsOnlyPunctuation(s string) bool {
	for _, r := range s {
		if !unicode.IsPunct(r) && !unicode.IsSpace(r) && !unicode.IsSymbol(r) {
			return false
		}
	}
	return true
}

// ShouldCountAsNotFound returns true if text is meaningful Japanese content
// that should have had a translation (not just punctuation/symbols)
func ShouldCountAsNotFound(s string) bool {
	if s == "" {
		return false
	}
	if IsOnlyPunctuation(s) {
		return false
	}
	return ContainsJapanese(s)
}

// PatchStats tracks translation application statistics during patching.
// Use StartTracking/StopTracking to enable stats collection via DictLookup.
type PatchStats struct {
	Applied  int                    // number of successful lookups
	NotFound int                    // number of non-empty dialogue texts with no translation
	UsedKeys map[string]struct{}    // unique dictionary keys that were matched
}

var activeStats *PatchStats

// StartTracking begins collecting patch statistics. Returns the stats object
// that will be populated during subsequent DictLookup calls.
func StartTracking() *PatchStats {
	activeStats = &PatchStats{
		UsedKeys: make(map[string]struct{}),
	}
	return activeStats
}

// StopTracking stops collecting patch statistics and returns the final stats.
func StopTracking() *PatchStats {
	stats := activeStats
	activeStats = nil
	return stats
}

// DictLookup looks up a translation in the dictionary, using typed keys when keyMode is "typed",
// falling back to legacy untyped keys for backward compatibility.
func DictLookup(dictionary map[string]string, keyMode string, entryType string, text string) (string, bool) {
	if keyMode == "typed" {
		key := GetTypedTranslationKey(entryType, text)
		if translation, ok := dictionary[key]; ok {
			if activeStats != nil {
				activeStats.Applied++
				activeStats.UsedKeys[key] = struct{}{}
			}
			return translation, true
		}
		// Fallback to untyped key for backward compat within mixed patches
		key = GetTranslationKey(text)
		if translation, ok := dictionary[key]; ok {
			if activeStats != nil {
				activeStats.Applied++
				activeStats.UsedKeys[key] = struct{}{}
			}
			return translation, true
		}
		if activeStats != nil && entryType == "dialogue" && ShouldCountAsNotFound(text) {
			activeStats.NotFound++
		}
		return "", false
	}
	// Legacy mode: untyped keys
	key := GetTranslationKey(text)
	if translation, ok := dictionary[key]; ok {
		if activeStats != nil {
			activeStats.Applied++
			activeStats.UsedKeys[key] = struct{}{}
		}
		return translation, true
	}
	if activeStats != nil && entryType == "dialogue" && ShouldCountAsNotFound(text) {
		activeStats.NotFound++
	}
	return "", false
}
