
down:
	@echo "Stopping and clear all"
	docker-compose down
	rmdir /s tmp

#i dont like to use docker compose with the api. Soo i create a env file with the docker compose config postgre for now.
up: 
	@echo "$(OK_COLOR)==>  Running in developer mode using docker-compose and air$(NO_COLOR)"
	@docker-compose up -d
	@sleep 3 && \
		echo PG_URI="postgres://test:test@`docker-compose port postgres 5432`/test?sslmode=disable" > .env