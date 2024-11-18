package main

import (
	"fmt"
	"image"
	"os"

	"gocv.io/x/gocv"
)

func cropper() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go [input image path] [output image path]")
		return
	}

	inputImagePath := os.Args[1]
	outputImagePath := os.Args[2]

	img := gocv.IMRead(inputImagePath, gocv.IMReadColor)
	if img.Empty() {
		fmt.Printf("Error reading image from: %v\n", inputImagePath)
		return
	}
	defer img.Close()

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	binary := gocv.NewMat()
	defer binary.Close()
	gocv.Threshold(gray, &binary, 250, 255, gocv.ThresholdBinaryInv)

	contours := gocv.FindContours(binary, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	var largestContour []image.Point
	maxArea := 0.0
	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		area := gocv.ContourArea(contour)
		if area > maxArea {
			maxArea = area
			largestContour = contour
		}
	}

	if len(largestContour) == 0 {
		fmt.Println("No document found in the image")
		return
	}

	rect := gocv.BoundingRect(largestContour)

	cropped := img.Region(rect)

	gocv.IMWrite(outputImagePath, cropped)

	fmt.Printf("Cropped image saved to: %v\n", outputImagePath)
}
