.PHONY: build
build:
	env GOOS=linux GOARCH=amd64 go build -v ./cmd/apiserver

.DEFAULT_GOAL := build