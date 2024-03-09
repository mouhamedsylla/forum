build:
	cd ./App && GOOS=js GOARCH=wasm go build -o ./internal/assets/main.wasm

start: build
	cd ./App/server && go run .

rm: 
	rm app

app: build
	cd ./cmd && go build -o ../app
	chmod u+x app