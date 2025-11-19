package processing

import "image"

func CalulateLightValue(img *image.Gray) float32 {
	bounds := img.Bounds()
	sizeX, sizeY := bounds.Max.X, bounds.Max.Y

	var light int64 = 0

	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			light += int64(img.GrayAt(x, y).Y)
		}
	}
	return float32(light) / (float32(sizeX) * float32(sizeY) * 255.0)
}
