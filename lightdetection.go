// Package lightdetection: Used to process images and build as wasm for react native application
package lightdetection

import "image"

type DetectedArea struct {
	Area        int
	BoundingBox image.Rectangle
	Centroid    image.Point
}
