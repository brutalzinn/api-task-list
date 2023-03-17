dev:
	docker-compose up -d
dev-build:
	docker-compose up --build

stop:
	@echo "Stopping and clear all"
	docker-compose down
	rmdir /s tmp