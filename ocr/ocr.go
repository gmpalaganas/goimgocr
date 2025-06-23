// Package ocr provides the functionilty to goimgocr to extract text from images
// using Google's Tesseract OCR engine.
package ocr

import (
	"errors"
	"genesis/goimgocr/ocr/imageprocessing"
	"genesis/goimgocr/ocr/languageutils"

	"github.com/otiai10/gosseract/v2"
)

const (
	tessdataDir     = "/usr/share/tessdata" // Directory where Tesseract language data files are stored
	targetPixelArea = 1000.0                // Target pixel area for image preprocessing
)

// ExtractTextFromImage takes the path to an image file,
// and the languages to extract.
//
// Returns the extracted text as a string.
func ExtractTextFromImage(imagePath string, languages ...string) (string, error) {
	// Check if provided language is available
	for _, lang := range languages {
		languageExists, _ := languageutils.CheckLanguageExists(tessdataDir, lang)
		if !languageExists {
			langErr := errors.New("language " + lang + " does not exist in tessdata directory")
			return "", langErr
		}
	}

	ocrClient := gosseract.NewClient()
	defer ocrClient.Close()

	// Set the languages to be used for OCR.
	err := ocrClient.SetLanguage(languages...)
	if err != nil {
		return "", err
	}

	imageBytes, err := imageprocessing.PreprocessImage(imagePath, targetPixelArea)
	if err != nil {
		return "", err
	}

	ocrClient.SetImageFromBytes(imageBytes)

	return ocrClient.Text()
}
