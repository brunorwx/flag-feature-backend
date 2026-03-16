.PHONY: help build run 

help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Build and run the application"

build:
	@echo "Building Feature Flag API..."
	go build -o flagapi ./cmd/main.go

run: build
	@echo "Running Feature Flag API..."
	./flagapi


