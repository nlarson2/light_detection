//go:build js && wasm

// Package wasm used to compile wasm golang code
package main

import (
	"encoding/json"
	"syscall/js"

	"light_detection/processing"
)

func jsDetectFlashlight(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf(map[string]any{"error": "missing base64"})
	}

	images := args[0]
	// b64 := args[0].String()

	// Defaults tuned for phone flashlight in a dark-ish scene
	pthresh := 0.995 // keep top ~0.5% brightest pixels

	if len(args) >= 2 {
		pthresh = args[1].Float()
	}

	var outVals []bool
	n := images.Get("length").Int()
	for i := 0; i < n; i++ {
		b64 := images.Index(i).String()
		img, err := processing.DecodeBase64ToJpeg(b64)
		if err != nil {
			return nil
		}

		gray := processing.ImageToGray(img)
		processing.ThresholdOfGrayImage(gray, float32(pthresh))
		_, detected := processing.KeepLargestArea(gray, 15)
		err = processing.EncodeJpegToFile(gray, "./images/test")
		outVals = append(outVals, detected)
		// outVals = append(outVals, processing.CalulateLightValue(gray))
	}

	res := map[string]any{
		"light_detected": outVals,
	}
	b, _ := json.Marshal(res)
	var out map[string]any
	_ = json.Unmarshal(b, &out)
	return js.ValueOf(out)
}

func main() {
	js.Global().Set("goDetectFlashlightJpeg", js.FuncOf(jsDetectFlashlight))
	select {}
}
