.PHONY: help build run dev tidy test

help:
	@echo "Commands:"
	@echo "  dev      Run development server"
	@echo "  build    Build binary"
	@echo "  tidy     Tidy dependencies"
	@echo "  test     Run tests"

dev:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

tidy:
	go mod tidy

test:
	go test -v ./...