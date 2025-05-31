.PHONY: build run up down logs restart

# Build the Docker image
build:
	docker-compose up --build

# Run the container (standalone)
run:
	docker-compose -f ./docker-compose.yml up

# Stop the stack
down:
	docker-compose down

# Show logs from Docker Compose
logs:
	docker-compose logs -f --tail=100

# Restart the stack
restart: down up
