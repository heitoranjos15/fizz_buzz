# Makefile for Go API with Gin
# Variables
APP_NAME = my-go-api
DOCKER_IMAGE = $(APP_NAME):latest
PORT = 8080

# Default target
.PHONY: all
all: build

# Build the Go application locally
.PHONY: build
build:
	go build -o $(APP_NAME) ./main.go

# Run the Go application locally
.PHONY: run
run:
	go run ./main.go

# Run tests
.PHONY: test
test:
	go test ./... -v

# Format the Go code
.PHONY: fmt
fmt:
	go fmt ./...

# Build the Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Run the Docker container
.PHONY: docker-run
docker-run:
	docker run -p $(PORT):$(PORT) $(DOCKER_IMAGE)

# Stop and remove the Docker container
.PHONY: docker-stop
docker-stop:
	docker stop $(shell docker ps -q --filter ancestor=$(DOCKER_IMAGE)) || true
	docker rm $(shell docker ps -a -q --filter ancestor=$(DOCKER_IMAGE)) || true

# Clean up local build artifacts
.PHONY: clean
clean:
	rm -f $(APP_NAME)
	go clean

# Clean up Docker images
.PHONY: docker-clean
docker-clean:
	docker rmi $(DOCKER_IMAGE) || true

# Run tests and build the Docker image
.PHONY: ci
ci: test docker-build

# Help command to display available targets
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build        - Build the Go application locally"
	@echo "  make run          - Run the Go application locally"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format Go code"
	@echo "  make docker-build - Build the Docker image"
	@echo "  make docker-run   - Run the Docker container"
	@echo "  make docker-stop  - Stop and remove the Docker container"
	@echo "  make clean        - Clean up local build artifacts"
	@echo "  make docker-clean - Clean up Docker images"
	@echo "  make ci           - Run tests and build Docker image"
	@echo "  make help         - Show this help message"
