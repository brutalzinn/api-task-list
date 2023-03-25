user := root
database := test

dev:
	docker-compose up -d
	@echo "docker development setup started."
	
dev-build:
	docker-compose up --build
	@echo "docker compose image rebuilded."

stop:
	@echo "Stopping and clear all"
	docker-compose down
	@echo "Docker compose Stopped"

swagger:
	docker-compose exec -i app swag init
	@echo "Swagger doc generated"

create-db:
	docker-compose exec -it postgres psql -U ${user} -d ${database} -f /scripts/init.sql
	@echo "database created"

delete-data-db:
	docker-compose exec -it postgres psql -U ${user} -d ${database} -f /scripts/clear_data.sql
	@echo "database clearned"

insert-data-db:
	docker-compose exec -it postgres psql -U ${user} -d ${database} -f /scripts/insert_data.sql
	@echo "Initial data inserted"

restart-db: delete-data-db insert-data-db