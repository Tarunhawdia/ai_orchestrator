.PHONY: all build run test clean docker-build docker-run

# Define service paths
ORCHESTRATOR_PATH = ./services/orchestrator

# Default target
all: build

# Build all Go services
build:
	@echo "Building all Go services..."
	go build -o $(ORCHESTRATOR_PATH)/main $(ORCHESTRATOR_PATH)
	@echo "Build complete."

# Run the orchestrator service locally
run: build
	@echo "Running Orchestrator Service..."
	$(ORCHESTRATOR_PATH)/main

# Run the orchestrator service locally with a specific port
run-port: build
	@echo "Running Orchestrator Service on port 8081..."
	PORT=8081 $(ORCHESTRATOR_PATH)/main

# Run all tests
test:
	@echo "Running all tests..."
	go test ./... -v

# Clean up built binaries
clean:
	@echo "Cleaning up..."
	go clean -i ./...
	rm -f $(ORCHESTRATOR_PATH)/main
	@echo "Clean complete."

# Build Docker image for orchestrator
docker-build:
	@echo "Building Docker image for Orchestrator..."
	docker build -t decentralized-ai-orchestrator-orchestrator -f $(ORCHESTRATOR_PATH)/Dockerfile .
	@echo "Docker image build complete."

# Run Docker container for orchestrator
docker-run: docker-build
	@echo "Running Orchestrator Docker container..."
	docker run --rm -p 8080:8080 --name orchestrator-container decentralized-ai-orchestrator-orchestrator
	@echo "To stop: docker stop orchestrator-container"

.PHONY: setup-env
# Setup .env file from .env.example
setup-env:
	@echo "Setting up .env file from .env.example..."
	cp .env.example .env
	@echo ".env file created. Please edit it with your local configurations."