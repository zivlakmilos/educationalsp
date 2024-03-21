.PHONY: build test

build:
	go build -o build/educationalsp ./cmd/educationalsp

test:
	go test ./...

run: build
	go run ./cmd/educationalsp

all: run
