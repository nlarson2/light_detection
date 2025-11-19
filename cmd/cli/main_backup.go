// Package cli used to test code before compiling to wasm
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"

	IMG "light_detection/image"
)

func getStructBytes(s interface{}) []byte {
	// Get the size of the struct in bytes
	size := unsafe.Sizeof(s)

	// Get a pointer to the struct
	ptr := unsafe.Pointer(&s)

	// Create a slice header pointing to the struct's memory
	var sliceHeader reflect.SliceHeader
	sliceHeader.Data = uintptr(ptr)
	sliceHeader.Len = int(size)
	sliceHeader.Cap = int(size)

	// Convert the slice header to a byte slice
	return *(*[]byte)(unsafe.Pointer(&sliceHeader))
}

func backup_main() {
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
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		for i := range 50 {
			fmt.Println(content[i])
		}
		return
		image, err := IMG.DecodeJpegFromBytes(content)
		IMG.ConvertToGrayScale(image)
		fmt.Println("Pixels: ", len(image.Pixels), " VALUES: ", len(image.Pixels)*3)

		// fmt.Println(content)
		b := make([]byte, len(image.Pixels)*3)

		for _, pixel := range image.Pixels {
			b = append(b, byte(pixel.R))
			b = append(b, byte(pixel.G))
			b = append(b, byte(pixel.B))
		}

		fmt.Println("Bytes: ", len(b))
		_ = os.WriteFile(write_path+f, b, 0o644)

	}
}
