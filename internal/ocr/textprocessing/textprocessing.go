// Package textprocessing provides utilities for processing the output of the OCR
package textprocessing

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

type textLine struct {
	index int
	text  string
}

// TrimJapanseText trims the surrounding spaces from Japanese characters
// in the given text.
func TrimJapanseText(text string) (string, error) {
	lines := strings.SplitAfter(text, "\n")
	newLines := make([]string, len(lines))

	var lineProcessingGroup sync.WaitGroup

	linesChan := make(chan textLine, len(lines))

	for i, line := range lines {
		// Launch a goroutine for each line to process it concurrently
		lineProcessingGroup.Add(1)
		go func() {
			defer lineProcessingGroup.Done()
			if hasJapaneseCharacters(line) {
				trimPattern := `\s+([\p{Hiragana}\p{Katakana}\p{Han}])\s+`
				surroundingSpaceRegexp := regexp.MustCompile(trimPattern)
				newLine := surroundingSpaceRegexp.ReplaceAllString(line, "$1")
				newLine = strings.TrimSpace(newLine)
				linesChan <- textLine{index: i, text: newLine}
			} else {
				linesChan <- textLine{index: i, text: line}
			}
		}()
	}

	// Wait for all line processing goroutines to finish
	lineProcessingGroup.Wait()
	close(linesChan)

	for newLine := range linesChan {
		fmt.Println("Processing line:", newLine.index, "Text:", newLine.text)
		newLines[newLine.index] = newLine.text
	}
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
