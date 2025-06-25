// Package imageprocessing provides utility functions for
// preprocessing images before passing them to the OCR engine.
package imageprocessing

import (
	"image"
	"math"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/cuda"
)

type IMProcessingMode int

const (
	NoCUDA IMProcessingMode = iota
	CUDA
)

// PreprocessImage takes the path to an image file and a target pixel area
// Processes the image for better OCR results
// NOTE: This preprocessing is only based on what worked for me for JPN and ENG languages in the past
func PreprocessImage(imagePath string, targetPixelArea float64, mode IMProcessingMode) ([]byte, error) {
	preprocessedImageMat := gocv.IMRead(imagePath, gocv.IMReadColor)
	defer preprocessedImageMat.Close()

	var outputBytes []byte
	var err error

	switch mode {
	case CUDA:
		outputBytes, err = processImageCuda(preprocessedImageMat, targetPixelArea)
	default:
		outputBytes, err = processImage(preprocessedImageMat, targetPixelArea)
	}

	if err != nil {
		return nil, err
	}

	return outputBytes, nil
}

// Calculates the scaling factor needed for input image to reach target pixel area
func computScalingFactor(imageSize image.Point, targetPixelArea float64) (float64, error) {
	// Convert the input Mat to an image and calculate its pixel area

	imgCurDensity := imageSize.X * imageSize.Y

	scalingFactor := targetPixelArea / float64(imgCurDensity)

	return math.Sqrt(scalingFactor), nil
}

func processImage(imageMat gocv.Mat, targetPixelArea float64) ([]byte, error) {
	// Convert image to grayscale
	err := gocv.CvtColor(imageMat, &imageMat, gocv.ColorBGRToGray)
	if err != nil {
		return nil, err
	}

	inputImage, err := imageMat.ToImage()
	if err != nil {
		return nil, err
	}

	imageSize := inputImage.Bounds().Size()

	// If image pixel areaa is less than target pixel area then scale the image
	// to meet the target pixel area
	if float64(imageSize.X*imageSize.Y) < targetPixelArea {
		scalingFactor, err := computScalingFactor(imageSize, targetPixelArea)
		if err != nil {
			return nil, err
		}

		// Scale the image to reach the target pixel area
		err = gocv.Resize(
			imageMat,
			&imageMat,
			image.Pt(0, 0),
			scalingFactor, scalingFactor,
			gocv.InterpolationLinear)
		if err != nil {
			return nil, err
		}
	}

	// Apply thresholding to binarize the image
	gocv.Threshold(
		imageMat,
		&imageMat,
		0,
		255,
		gocv.ThresholdBinary|gocv.ThresholdOtsu)

	// Convert the processed Mat to gocv.NativeByteBuffer
	imageBytes, err := gocv.IMEncode(".png", imageMat)
	if err != nil {
		return nil, err
	}

	return imageBytes.GetBytes(), nil
}

func processImageCuda(imageMat gocv.Mat, targetPixelArea float64) ([]byte, error) {
	imageCudaMat := cuda.NewGpuMatFromMat(imageMat)
	defer imageCudaMat.Close()

	outputMat := imageMat.Clone()
	defer outputMat.Close()

	err := cuda.CvtColor(imageCudaMat, &imageCudaMat, gocv.ColorBGRToGray)
	if err != nil {
		return nil, err
	}

	inputImage, err := imageMat.ToImage()
	if err != nil {
		return nil, err
	}

	imageSize := inputImage.Bounds().Size()

	// If image pixel areaa is less than target pixel area then scale the image
	// to meet the target pixel area
	if float64(imageSize.X*imageSize.Y) < targetPixelArea {
		scalingFactor, err := computScalingFactor(imageSize, targetPixelArea)
		if err != nil {
			return nil, err
		}

		// Scale the image to reach the target pixel area
		err = cuda.Resize(
			imageCudaMat,
			&imageCudaMat,
			image.Pt(0, 0),
			scalingFactor, scalingFactor,
			cuda.InterpolationLinear)
		if err != nil {
			return nil, err
		}
	}

	// Apply thresholding to binarize the image
	cuda.Threshold(
		imageCudaMat,
		&imageCudaMat,
		0,
		255,
		gocv.ThresholdBinary|gocv.ThresholdOtsu)

	imageCudaMat.Download(&outputMat)

	// Convert the processed Mat to cuda.NativeByteBuffer
	imageBytes, err := gocv.IMEncode(".png", outputMat)
	if err != nil {
		return nil, err
	}

	return imageBytes.GetBytes(), nil
}
