dev: 
	@docker compose up
run:
	@go run ./cmd/api

stop all:
	@docker compose down	
live:
	@air

asynq web:
	@docker run --rm --name asynqmon -p  8080:8080 hibiken/asynqmon --redis-addr=host.docker.internal:6379

mailer test:
	openssl s_client -starttls smtp -connect localhost:1025