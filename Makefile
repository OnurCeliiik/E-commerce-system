make up:
	@echo "Starting services..."
	docker compose up -d
	@echo "Services started successfully"

make full up:
	@echo "Building and starting services..."
	docker compose up -d --build
	@echo "Services started successfully"

make down:
	@echo "Stopping services..."
	docker compose down
	@echo "Services stopped successfully"

make full down:
	@echo "Stopping and removing services..."
	docker compose down -v
	@echo "Services stopped and removed successfully"
