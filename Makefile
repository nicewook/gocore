.PHONY: test
test:
	@echo "Running tests..."
	@echo "Make mock first..."
	mockery
	@echo "Running tests..."
	go test ./... -v

# Makefile for Docker Compose Operations
.PHONY: up down down-v
up: # Docker Compose Up
	docker compose up -d

down: # Docker Compose Down (without volume removal)
	docker compose down

down-v: # Docker Compose Down with Volume Removal
	docker compose down -v
