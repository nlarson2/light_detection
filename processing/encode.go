package processing

import (
	"image"
	"image/jpeg"
	"os"
)

func EncodeJpegToFile(img image.Image, path string) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return jpeg.Encode(outputFile, img, nil)
}
