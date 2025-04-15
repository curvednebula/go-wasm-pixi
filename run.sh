set -e
GOOS=js GOARCH=wasm go build -o main.wasm
basic-http-server .