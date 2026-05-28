set shell := ["bash", "-c"]

default:
	@just --list

build:
    go build -o bin/lumingrate -ldflags="-s -w" main.go

run:
    go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run

alias fmt := format
format:
	go fmt ./...
	golangci-lint fmt