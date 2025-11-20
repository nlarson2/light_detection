// Package cli used to test code before compiling to wasm
package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"light_detection/processing"
)

func main() {
	path := "./test_images/preprocessed/"
	write_path := "./test_images/postprocessed/"
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
		} else {
			files = append(files, name)
		}
	}
	// If you want full paths:
	for _, f := range files {
		filePath := filepath.Join(path, f)
		jpegFile, _ := os.Open(filePath)

		defer jpegFile.Close()

		img, _ := jpeg.Decode(jpegFile)
		gray := processing.ImageToGray(img)
		if err := processing.EncodeJpegToFile(gray, write_path+f); err != nil {
			log.Fatal("Failed to encode jpeg file: ", err)
		}
		_ = processing.ThresholdOfGrayImage(gray, 0.995)
		minArea := 0.0005 * float64(gray.Bounds().Max.X*gray.Bounds().Max.Y)
		box, detected := processing.KeepLargestArea(gray, int(minArea))
		outputFile, _ := os.Create(write_path + f)
		defer outputFile.Close()

		_ = jpeg.Encode(outputFile, gray, nil)

		fmt.Println("File: ", f, "   Light Detected: ", processing.CalulateLightValue(gray), "  Detected: ", detected, "  Box: ", box, " Size: ", gray.Bounds(), " MINAREA: ", minArea)

	}
}
