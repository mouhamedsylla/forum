cd ./App
GOOS=js GOARCH=wasm go build -o ./internal/assets/main.wasm
cd ../cmd
go run .