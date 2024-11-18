package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// loadImage loads the image file and decodes it based on its format.
func loadImage(filename string) (image.Image, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

// getExifOrientation extracts the EXIF orientation value from the image file.
func getExifOrientation(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 1, err // Assume default orientation (1) if error occurs
	}
	defer file.Close()

	exifData, err := exif.Decode(file)
	if err != nil {
		return 1, err // Assume default orientation (1) if no EXIF data
	}

	orientationTag, err := exifData.Get(exif.Orientation)
	if err != nil {
		return 1, err // Assume default orientation (1) if no orientation tag
	}

	orientation, err := orientationTag.Int(0)
	if err != nil {
		return 1, err // Assume default orientation (1) if invalid orientation
	}

	return orientation, nil
}

// autoRotate rotates the image based on the EXIF orientation.
func autoRotate(img image.Image, orientation int) image.Image {
	switch orientation {
	case 3:
		// 180 degrees
		return imaging.Rotate180(img)
	case 6:
		// 90 degrees clockwise
		return imaging.Rotate270(img)
	case 8:
		// 90 degrees counter-clockwise
		return imaging.Rotate90(img)
	default:
		// No rotation needed
		return img
	}
}

// saveImage saves the rotated image back to the file.
func saveImage(img image.Image, format string, outputFile string) error {
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(out, img, nil)
	case "png":
		err = png.Encode(out, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	return err
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <image-file>\n", filepath.Base(os.Args[0]))
	}

	// Load the image
	imageFile := os.Args[1]
	img, format, err := loadImage(imageFile)
	if err != nil {
		log.Fatalf("Error loading image: %v\n", err)
	}

	// Get the EXIF orientation
	orientation, err := getExifOrientation(imageFile)
	if err != nil {
		log.Printf("Warning: Could not get EXIF orientation, assuming default: %v\n", err)
	}

	// Auto-rotate the image based on EXIF orientation
	rotatedImg := autoRotate(img, orientation)

	// Save the corrected image
	outputFile := "corrected_" + filepath.Base(imageFile)
	err = saveImage(rotatedImg, format, outputFile)
	if err != nil {
		log.Fatalf("Error saving image: %v\n", err)
	}

	fmt.Printf("Image successfully auto-rotated and saved as %s\n", outputFile)
}
