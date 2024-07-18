LOCAL_BIN ?= ./.env

.PHONY: demo.build
demo.build:
	cd demo && GOOS=js GOARCH=wasm go build -o main.wasm

.PHONY: demo.serve
demo.serve:
	cd demo && go run server/main.go

.PHONY: install.golangci
install.golangci:
	mkdir -p $(LOCAL_BIN) && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN) v1.56.2

.PHONY: lint
lint:
	./.env/golangci-lint run --skip-dirs=/demo/*

.PHONY: example.simple
example.simple:
	SOLVE_ALL=true go run example/simple/main.go sample/words.json