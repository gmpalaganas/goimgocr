// Package imageprocessing provides utility functions for
// preprocessing images before passing them to the OCR engine.
package imageprocessing

import (
	"fmt"
	"image"
	"math"

	"gocv.io/x/gocv"
)

// PreprocessImage takes the path to an image file and a target pixel area
// Processes the image for better OCR results
// NOTE: This preprocessing is only based on what worked for me for JPN and ENG languages in the past
func PreprocessImage(imagePath string, targetPixelArea float64) ([]byte, error) {
	preprocessedImageMat := gocv.IMRead(imagePath, gocv.IMReadColor)

	// Convert image to grayscale
	err := gocv.CvtColor(preprocessedImageMat, &preprocessedImageMat, gocv.ColorBGRToGray)
	if err != nil {
		return nil, err
	}

	preprocessedImage, err := preprocessedImageMat.ToImage()
	if err != nil {
		return nil, err
	}

	preProcessedImageSize := preprocessedImage.Bounds().Size()

	// If image pixel areaa is less than target pixel area then scale the image
	// to meet the target pixel area
	fmt.Println("Preprocessed image size:", preProcessedImageSize.X*preProcessedImageSize.Y)
	if float64(preProcessedImageSize.X*preProcessedImageSize.Y) < targetPixelArea {
		scalingFactor, err := computScalingFactor(preProcessedImageSize, targetPixelArea)
		if err != nil {
			return nil, err
		}

		// Scale the image to reach the target pixel area
		err = gocv.Resize(
			preprocessedImageMat,
			&preprocessedImageMat,
			image.Pt(0, 0),
			scalingFactor, scalingFactor,
			gocv.InterpolationLinear)
		if err != nil {
			return nil, err
		}
	}

	// Apply thresholding to binarize the image
	gocv.Threshold(
		preprocessedImageMat,
		&preprocessedImageMat,
		0,
		255,
		gocv.ThresholdBinary|gocv.ThresholdOtsu)

	// Convert the processed Mat to gocv.NativeByteBuffer
	preprocessedImageBytes, err := gocv.IMEncode(".png", preprocessedImageMat)
	if err != nil {
		return nil, err
	}

	return preprocessedImageBytes.GetBytes(), nil
}

// Calculates the scaling factor needed for input image to reach target pixel area
func computScalingFactor(imageSize image.Point, targetPixelArea float64) (float64, error) {
	// Convert the input Mat to an image and calculate its pixel area

	imgCurDensity := imageSize.X * imageSize.Y

	scalingFactor := targetPixelArea / float64(imgCurDensity)

	return math.Sqrt(scalingFactor), nil
}
