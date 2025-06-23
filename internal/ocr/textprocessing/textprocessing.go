// Package textprocessing provides utilities for processing the output of the OCR
package textprocessing

import (
	"regexp"
	"strings"
	"unicode"
)

func TrimJapanseText(text string) (string, error) {
	japaneseCharPattern := `[\p{Hiragana}\p{Katakana}\p{Han}]`

	spaceBeforeJpRegexp := regexp.MustCompile(`\s+(` + japaneseCharPattern + `)`)
	spaceAfterJpRegexp := regexp.MustCompile(`(` + japaneseCharPattern + `)\s+`)

	outputText := ""

	lines := strings.SplitAfterSeq(text, "\n")

	for line := range lines {
		if hasJapaneseCharacters(line) {
			lineText := spaceBeforeJpRegexp.ReplaceAllString(line, "$1")
			outputText += spaceAfterJpRegexp.ReplaceAllString(lineText, "$1") + "\n"
		} else {
			// Skip lines that do not contain Japanese characters
			outputText += strings.TrimSpace(line) + "\n"
		}
	}

	return outputText, nil
}

// hasJapaneseCharacters checks if the given text contains any Japanese characters
func hasJapaneseCharacters(text string) bool {
	for _, char := range text {
		if unicode.Is(unicode.Hiragana, char) ||
			unicode.Is(unicode.Katakana, char) ||
			unicode.Is(unicode.Unified_Ideograph, char) {

			return true
		}
	}
	return false
}
