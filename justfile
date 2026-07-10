set shell := ["bash", "-c"]

mod web

default:
	@just --list

cibuild: build-bin build-wasm

build: build-bin build-wasm web::build

build-bin:
    go build -o bin/lumingrate -ldflags="-s -w" main.go

[env("GOOS", "js")]
[env("GOARCH", "wasm")]
build-wasm:
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