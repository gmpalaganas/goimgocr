// Package languageutils provides utility functions for handling
// language options for gosseract
package languageutils

import (
	"io/fs"
	"os"
	"slices"
	"strings"
)

const (
	trainedDataExt = ".traineddata"
)

// GetAvailableLanguages gets the available language for tesseract
// from the provided tessdata dir path
func GetAvailableLanguages(tessDataDir string) ([]string, error) {
	languages := []string{}

	tessDataDirFS := os.DirFS(tessDataDir)
	tessDataFiles, err := fs.Glob(tessDataDirFS, "*.traineddata")
	if err != nil {
		return nil, err
	}

	for _, tessDataFile := range tessDataFiles {
		language := strings.Replace(tessDataFile, trainedDataExt, "", 1)
		language = strings.Replace(language, tessDataDir, "", 1)
		languages = append(languages, language)
	}

	return languages, nil
}

// CheckLanguageExists checks if the provided language exists in the
// provided tessdata dir path
func CheckLanguageExists(tessDataDir, language string) (bool, error) {
	languages, err := GetAvailableLanguages(tessDataDir)

	return slices.Contains(languages, language), err
}
