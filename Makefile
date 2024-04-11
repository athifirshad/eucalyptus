build:
	@docker compose up -build
dev: 
	@docker compose  -f "docker-compose.yml" up -d --build redis 
run:
	@go run ./cmd/api

stop:
	@docker compose down	
live:
	@air

asynq web:
	@docker run --rm --name asynqmon -p  8080:8080 hibiken/asynqmon --redis-addr=host.docker.internal:6379

mailer test:
	openssl s_client -starttls smtp -connect localhost:1025

psql:
	docker exec -it db psql -U root eucalyptus

createdb:
	docker exec -it db createdb --user=root --owner=root eucalyptus
	
package:
	@docker build -t ghcr.io/athifirshad/eucalyptus .

push:
	@docker push ghcr.io/athifirshad/eucalyptus