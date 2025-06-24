// Package goimgocr takes an image and extract the text from it using
// Google's Tesseract OCR engine.
// Currently only supports English and Japanese languages.
package main

import (
	"flag"
	"fmt"
	"genesis/goimgocr/internal/ocr"
	"log"
	"os"
	"strings"
	"time"
)

// NOTE: Temp default values
// TODO: Make main read from a config file or command line arguments
const (
	tessdataDirDefault     = "/usr/share/tessdata"    // Directory where Tesseract language data files are stored
	targetPixelAreaDefault = 500000.0                 // Target pixel area for image preprocessing
	languagesDefault       = "jpn+eng"                // Languages to be used for OCR
	testImagePathDefault   = "./testimg/testimg1.png" // Path to the image file to be processed
)

func main() {
	args := os.Args

	if len(args) <= 1 {
		log.Fatal("Please provide the path to the image file as an argument.")
	}

	// Set up command line flags
	var tessDataDir string
	var targetPixelArea float64
	var languages string
	var debugMode bool

	var startTime time.Time

	flag.StringVar(
		&tessDataDir,
		"tessdata",
		tessdataDirDefault,
		"Directory where Tesseract language data files are stored")
	flag.Float64Var(
		&targetPixelArea,
		"target-pixel-area",
		targetPixelAreaDefault,
		"Target pixel area for image preprocessing")
	flag.StringVar(
		&languages,
		"languages",
		languagesDefault,
		"Languages to be used for OCR ('+'-separated)")
	flag.BoolVar(
		&debugMode,
		"debug",
		false,
		"Enable or disable debug mode")

	// Set Usage message
	usage := func() {
		fmt.Printf("Usage: %s [options] [image_path]\nOptions:\n", args[0])
		flag.PrintDefaults()
	}

	flag.Usage = usage

	flag.Parse()

	imagePath := args[len(args)-1] // Get the last argument as the image path

	// Check if the image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("Image file does not exist: %s", imagePath)
	}

	languagesList := strings.Split(languages, "+")

	config := ocr.OCRConfig{
		TessDataDir:     tessDataDir,
		TargetPixelArea: targetPixelArea,
		Languages:       languagesList,
	}

	if debugMode {
		fmt.Println("Running with the following configuration:")
		fmt.Println("Tesseract data Directory:", config.TessDataDir)
		fmt.Println("Target Pixel Area:", config.TargetPixelArea)
		fmt.Println("Languages:", strings.Join(config.Languages, ", "))

		startTime = time.Now()
	}

	text, err := ocr.ExtractTextFromImage(imagePath, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(text)

	if debugMode {
		fmt.Println("\nOCR finished in ", time.Since(startTime))
	}
}
