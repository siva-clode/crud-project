APP_NAME=crud-project
IMAGE_NAME=crud-project
PORT=8080
COMPOSE_FILE=docker-compose.yml

.PHONY: help build run stop clean docker-build docker-run docker-stop docker-clean logs docker-up docker-down docker-ps docker-restart

help:
	@echo "🚀 Available commands:"
	@echo ""
	@echo "  make build           - Build the Go binary locally"
	@echo "  make run             - Run the app locally"
	@echo "  make stop            - Stop the local app"
	@echo "  make clean           - Remove build artifacts"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make docker-run      - Run single container (API only)"
	@echo "  make docker-stop     - Stop API container"
	@echo "  make docker-clean    - Remove API image & container"
	@echo "  make logs            - View logs from API container"
	@echo ""
	@echo "🧩 Docker Compose Commands:"
	@echo "  make docker-up       - Start API + PostgreSQL stack"
	@echo "  make docker-down     - Stop and remove Compose stack"
	@echo "  make docker-ps       - List running Compose services"
	@echo "  make docker-restart  - Rebuild and restart Compose stack"

# -----------------------------
# Local Development
# -----------------------------

build:
	@echo "🏗️  Building Go binary..."
	go build -o bin/$(APP_NAME) ./cmd/api

run:
	@echo "🚀 Running app locally..."
	./bin/$(APP_NAME)

stop:
	@echo "🛑 Stopping app..."
	pkill -f "./bin/$(APP_NAME)" || true

clean:
	@echo "🧹 Cleaning build files..."
	rm -rf bin

# -----------------------------
# Single Container (API only)
# -----------------------------

docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(IMAGE_NAME) .

docker-run:
	@echo "🚀 Running Docker container on port $(PORT)..."
	docker run --env-file .env -p $(PORT):8080 --name $(APP_NAME)-container $(IMAGE_NAME)

docker-stop:
	@echo "🛑 Stopping Docker container..."
	docker stop $(APP_NAME)-container || true
	docker rm $(APP_NAME)-container || true

docker-clean: docker-stop
	@echo "🧹 Removing Docker image..."
	docker rmi $(IMAGE_NAME) || true

logs:
	docker logs -f $(APP_NAME)-container

# -----------------------------
# Docker Compose (API + DB)
# -----------------------------

docker-up:
	@echo "🐳 Starting full stack (API + DB)..."
	docker compose -f $(COMPOSE_FILE) up --build -d
	@echo "✅ Stack is up and running!"

docker-down:
	@echo "🛑 Stopping and removing stack..."
	docker compose -f $(COMPOSE_FILE) down -v
	@echo "🧹 All containers and volumes removed."

docker-ps:
	docker compose -f $(COMPOSE_FILE) ps

docker-restart:
	@echo "♻️  Restarting stack..."
	docker compose -f $(COMPOSE_FILE) down -v
	docker compose -f $(COMPOSE_FILE) up --build -d
	@echo "✅ Stack restarted successfully!"
