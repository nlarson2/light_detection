package processing

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
)

func DecodeBase64ToJpeg(data string) (image.Image, error) {
	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(bytes.NewReader(raw))
}
