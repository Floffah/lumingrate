set shell := ["bash", "-c"]

default:
	@just --list

build: build-web
    go build -o bin/lumingrate -ldflags="-s -w" main.go

[env("GOOS", "js")]
[env("GOARCH", "wasm")]
build-web:
	go build -o web/src/lib/lumingrate.wasm -ldflags="-s -w" cmd/web/web.go

run:
    go run main.go

runc CHAPTER:
	go run main.go --chapter {{CHAPTER}}

test:
	go test -v ./...

lint:
	golangci-lint run

alias fmt := format
format:
	go fmt ./...
	golangci-lint fmt