.PHONY: demo.build
demo.build:
	cd demo && GOOS=js GOARCH=wasm go build -o main.wasm

.PHONY: demo.serve
demo.serve:
	cd demo && go run server/main.go