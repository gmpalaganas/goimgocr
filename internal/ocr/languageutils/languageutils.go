// Package languageutils provides utility functions for handling
// language options for gosseract
package languageutils

import (
	"os/exec"
	"slices"
	"strings"
)

// GetAvailableLanguages gets the available language for tesseract
// from the provided tessdata dir path
func GetAvailableLanguages(tessdataDir string) ([]string, error) {
	// If the tessdataDir does not end with a slash, append one
	// Exccept when empty string to avoid pointing to root directory
	if tessdataDir != "" && tessdataDir[len(tessdataDir)-1] != '/' {
		tessdataDir += "/"
	}

	cmdName := "find"
	cmdArgs := []string{tessdataDir, "-maxdepth", "1", "-type", "f", "-name", "*.traineddata"}

	cmdResult, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		return nil, err
	}

	// Parses the output of the command to get the list of available languages
	cmdResultStr := strings.ReplaceAll(string(cmdResult), tessdataDir, "")
	cmdResultStr = strings.ReplaceAll(cmdResultStr, ".traineddata", "")
	languages := strings.Split(strings.TrimSpace(cmdResultStr), "\n")

	return languages, nil
}

// CheckLanguageExists checks if the provided language exists in the
// provided tessdata dir path
func CheckLanguageExists(tessdataDir, language string) (bool, error) {
	languages, err := GetAvailableLanguages(tessdataDir)

	return slices.Contains(languages, language), err
}
