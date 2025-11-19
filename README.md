# WASM tool for light detection
## Using go version go1.24.4 linux/amd64
## 3 forms of running
1. Processing images - ex: go run cmd/cli/main.go<br>
    Note: Images should be stored in ./test_images/preprocessed/.
2. Web Server Display Images - ex:   go run cmd/image-viewer/main.go -root ./test_images/postprocessed -addr :8080<br>
or use the proprocessed directory if you want to see those
3. Compile to wasm - ex: GOOS=js GOARCH=wasm go build -o ./main.wasm ./cmd/wasm/main.go
