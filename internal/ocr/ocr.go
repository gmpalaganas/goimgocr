// Package ocr provides the functionilty to goimgocr to extract text from images
// using Google's Tesseract OCR engine.
package ocr

import (
	"errors"
	"genesis/goimgocr/internal/ocr/imageprocessing"
	"genesis/goimgocr/internal/ocr/languageutils"
	"genesis/goimgocr/internal/ocr/textprocessing"

	"github.com/otiai10/gosseract/v2"
)

type OCRConfig struct {
	TessDataDir     string                           // Directory where Tesseract language data files are stored
	TargetPixelArea float64                          // Target pixel area for image preprocessing
	Languages       []string                         // Languages to be used for OCR
	ProcessingMode  imageprocessing.IMProcessingMode // Processing mode for image preprocessing
}

// ExtractTextFromImage takes the path to an image file,
// and the languages to extract.
//
// Returns the extracted text as a string.
func ExtractTextFromImage(imagePath string, ocrConfig OCRConfig) (string, error) {
	// Check if provided language is available
	for _, lang := range ocrConfig.Languages {
		languageExists, _ := languageutils.CheckLanguageExists(ocrConfig.TessDataDir, lang)
		if !languageExists {
			langErr := errors.New("language " + lang + " does not exist in tessdata directory")
			return "", langErr
		}
	}

	ocrClient := gosseract.NewClient()
	defer ocrClient.Close()

	// Set the languages to be used for OCR.
	err := ocrClient.SetLanguage(ocrConfig.Languages...)
	if err != nil {
		return "", err
	}
	imageBytes, err := imageprocessing.PreprocessImage(imagePath, ocrConfig.TargetPixelArea, ocrConfig.ProcessingMode)
	if err != nil {
		return "", err
	}

	ocrClient.SetImageFromBytes(imageBytes)

	extractedText, err := ocrClient.Text()
	if err != nil {
		return "", err
	}

	outputText, err := textprocessing.TrimJapanseText(extractedText)
	if err != nil {
		return "", err
	}

	return outputText, nil
}
