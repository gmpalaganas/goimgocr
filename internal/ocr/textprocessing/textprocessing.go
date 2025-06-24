// Package textprocessing provides utilities for processing the output of the OCR
package textprocessing

import (
	"regexp"
	"strings"
	"sync"
	"unicode"
)

// TrimJapanseText trims the surrounding spaces from Japanese characters
// in the given text.
func TrimJapanseText(text string) (string, error) {
	lines := strings.SplitAfter(text, "\n")
	newLines := make([]string, len(lines))

	var lineProcessingGroup sync.WaitGroup

	for i, line := range lines {
		// Launch a goroutine for each line to process it concurrently
		go func() {
			if hasJapaneseCharacters(line) {
				lineProcessingGroup.Add(1)

				defer lineProcessingGroup.Done()
				trimPattern := `\s+([\p{Hiragana}\p{Katakana}\p{Han}])\s+`
				surroundingSpaceRegexp := regexp.MustCompile(trimPattern)
				newLine := surroundingSpaceRegexp.ReplaceAllString(line, "$1")
				newLine = strings.TrimSpace(newLine)
				newLines[i] = newLine + "\n"
			} else {
				newLines[i] = line
			}
		}()
	}

	// Wait for all line processing goroutines to finish
	lineProcessingGroup.Wait()

	outputText := strings.Join(newLines, "")

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
