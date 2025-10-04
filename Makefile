DOCKER_IMAGE = $(APP_NAME):latest
PORT = 8080

up:
	docker-compose up --build

down:
	docker-compose down

build:
	docker build -t my-go-app .

run:
	docker-compose up -d

test:
	go test ./... -v

fmt:
	go fmt ./...

mongo-shell:
	docker-compose exec mongo mongosh


# Help command to display available targets
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build        - Build the Go application locally"
	@echo "  make run          - Run the Go application locally"
	@echo "  make up           - Start the Docker containers"
	@echo "  make down         - Stop the Docker containers"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format Go code"
	@echo "  make mongo-shell  - Access the MongoDB shell"
