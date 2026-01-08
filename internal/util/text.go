package util

import (
	"regexp"
	"strings"
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

// Wrap wraps text to a specified width, accounting for RPG Maker placeholders
func Wrap(text string, width int) string {
	if width <= 0 {
		width = 58
	}

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

// GetTranslationKey generates a normalized key for dictionary lookup
func GetTranslationKey(text string) string {
	return strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), " ", ""))
}




