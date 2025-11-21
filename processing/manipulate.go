package processing

import (
	"fmt"
	"image"
	"image/color"
	"math"

	lightdetection "light_detection"
)

func ImageToGray(src image.Image) *image.Gray {
	b := src.Bounds()
	dst := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bb, _ := src.At(x, y).RGBA()
			R := int(r >> 8)
			G := int(g >> 8)
			B := int(bb >> 8)
			y8 := (299*R + 587*G + 114*B + 500) / 1000
			dst.Pix[y*dst.Stride+x] = uint8(y8)
		}
	}
	return dst
}

func ThresholdOfGrayImage(img *image.Gray, percentThreshold float32) error {
	var minLightValue uint8 = 230
	if percentThreshold > 1 || percentThreshold < 0 {
		return fmt.Errorf("Incorect value for percentThreshold (0.0 - 1.0)")
	}
	minVal, maxVal := uint8(math.MaxUint8), uint8(0)
	// minVal, maxVal = uint8(0), uint8(math.MaxUint8)
	// fmt.Println(minVal, maxVal)
	bounds := img.Bounds()
	sizeX, sizeY := bounds.Max.X, bounds.Max.Y

	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			value := img.GrayAt(x, y).Y
			if value < minVal {
				minVal = value
			}
			if value > maxVal {
				maxVal = value
			}
		}
	}

	// buffer := uint8(5)
	diff := maxVal - minVal
	thresholdValue := uint8((float32(diff) * percentThreshold)) + minVal
	if thresholdValue >= 255 {
		thresholdValue = uint8(255)
	}

	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			value := img.GrayAt(x, y).Y
			if value < thresholdValue || value < minLightValue {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}

	return nil
}

func KeepLargestArea(img *image.Gray, minArea int, maxArea int) (lightdetection.DetectedArea, bool) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return lightdetection.DetectedArea{}, false
	}

	type compInfo struct {
		area       int
		sumX, sumY int
		minX, minY int
		maxX, maxY int
	}

	label := make([]int, w*h) // component ID per pixel, 0 = background
	visited := make([]bool, w*h)
	comps := []compInfo{{}} // index 0 unused so IDs start at 1

	idx := func(x, y int) int {
		return (y-b.Min.Y)*w + (x - b.Min.X)
	}
	neighbors := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	nextID := 1

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			i := idx(x, y)
			if visited[i] {
				continue
			}
			if img.GrayAt(x, y).Y == 0 {
				visited[i] = true
				continue
			}

			id := nextID
			nextID++

			comp := compInfo{
				minX: x, maxX: x,
				minY: y, maxY: y,
			}

			q := [][2]int{{x, y}}
			visited[i] = true
			label[i] = id

			for len(q) > 0 {
				px, py := q[0][0], q[0][1]
				q = q[1:]

				comp.area++
				comp.sumX += px
				comp.sumY += py

				if px < comp.minX {
					comp.minX = px
				}
				if px > comp.maxX {
					comp.maxX = px
				}
				if py < comp.minY {
					comp.minY = py
				}
				if py > comp.maxY {
					comp.maxY = py
				}

				for _, d := range neighbors {
					nx, ny := px+d[0], py+d[1]
					if nx < b.Min.X || nx >= b.Max.X || ny < b.Min.Y || ny >= b.Max.Y {
						continue
					}
					ni := idx(nx, ny)
					if visited[ni] {
						continue
					}
					if img.GrayAt(nx, ny).Y == 0 {
						visited[ni] = true
						continue
					}
					visited[ni] = true
					label[ni] = id
					q = append(q, [2]int{nx, ny})
				}
			}

			comps = append(comps, comp)
		}
	}

	// Find largest component by area with minArea
	bestID := 0
	bestArea := 0
	for id, c := range comps {
		if id == 0 {
			continue
		}
		if c.area >= minArea && c.area <= maxArea && c.area > bestArea {
			bestArea = c.area
			bestID = id
		}
	}
	if bestID == 0 {
		return lightdetection.DetectedArea{}, false
	}

	// Second pass: zero out all pixelVs not in the best component
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			i := idx(x, y)
			if label[i] != bestID {
				img.SetGray(x, y, color.Gray{Y: 0})
			} else {
				img.SetGray(x, y, color.Gray{Y: 255})
			}
		}
	}

	c := comps[bestID]
	blob := lightdetection.DetectedArea{
		Area:        bestArea,
		BoundingBox: image.Rect(c.minX, c.minY, c.maxX+1, c.maxY+1),
		Centroid: image.Point{
			X: c.sumX / c.area,
			Y: c.sumY / c.area,
		},
	}
	return blob, true
}
