// Package goimgocr takes an image and extract the text from it using
// Google's Tesseract OCR engine.
// Currently only supports English and Japanese languages.
package main

import (
	"fmt"
	"genesis/goimgocr/internal/ocr"
	"log"
	"strings"
)

// NOTE: Temp values
// TODO: Make main read from a config file or command line arguments
const (
	tessdataDir     = "/usr/share/tessdata" // Directory where Tesseract language data files are stored
	targetPixelArea = 500000.0              // Target pixel area for image preprocessing
	languages       = "jpn eng"             // Languages to be used for OCR
)

func main() {
	fmt.Println("Starting OCR process...")
	languagesList := strings.Split(languages, " ")

	config := ocr.OCRConfig{
		TessdataDir:     tessdataDir,
		TargetPixelArea: targetPixelArea,
		Languages:       languagesList,
	}

	text, err := ocr.ExtractTextFromImage("./testimg/testimg1.png", config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Extracted Text:\n", text)
}
