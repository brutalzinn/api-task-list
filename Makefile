
dev: 
	@echo "$(OK_COLOR)==> Running in developer mode using docker-compose and air$(NO_COLOR)"
	@docker-compose up -d
	@sleep 3 && \
		PG_URI="postgres://test:test@`docker-compose port postgres 5432`/test?sslmode=disable" \
		air