# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application (backend + frontend concurrently)
run:
	@echo "Starting backend and frontend..."
	@trap 'kill 0' INT; \
	go run cmd/api/main.go & \
	( cd client && npm run dev ) & \
	wait

# Run Docker containers
docker-run:
	@echo "Starting Docker containers..."
	@if command -v docker >/dev/null && docker compose version >/dev/null 2>&1; then \
		docker compose up --build; \
	else \
		if command -v docker-compose >/dev/null; then \
			docker-compose up --build; \
		else \
			echo "Neither Docker Compose V2 nor V1 found."; \
			exit 1; \
		fi; \
	fi

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker compose down

# Test the application
test:
	@echo "Running tests..."
	@go test ./... -v

# Integration Tests
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload (with Air)
watch:
	@trap 'kill 0' INT; \
	if command -v air >/dev/null; then \
		air & \
		( cd client && npm run dev ) & \
		wait; \
	else \
		read -p "Go's 'air' is not installed. Install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air & \
			go run cmd/api/main.go & \
			( cd client && npm run dev ) & \
			wait; \
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

.PHONY: all build run test clean watch docker-run docker-down itest
