run:
	@docker compose up -d
	@go run ./cmd/api
live:
	@air