// Package goimgocr takes an image and extract the text from it using
// Google's Tesseract OCR engine.
// Currently only supports English and Japanese languages.
package main

import (
	"fmt"
	"genesis/goimgocr/internal/ocr"
	"log"
)

func main() {
	fmt.Println("Starting OCR process...")
	text, err := ocr.ExtractTextFromImage("./testimg/testimg1.png", "jpn", "eng")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Extracted Text:\n", text)
}
