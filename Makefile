.PHONY: test
test:
	@echo "Running tests..."
	@echo "Make mock first..."
	mockery
	@echo "Running tests..."
	go test ./... -v
