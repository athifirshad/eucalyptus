run:
	@docker compose up -d
	@go run ./cmd/api
live:
	@air

asynq web:
	@docker run --rm --name asynqmon -p  8080:8080 hibiken/asynqmon --redis-addr=host.docker.internal:6379