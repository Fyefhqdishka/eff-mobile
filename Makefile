up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

run-tests:
	go test -v ./internal/handlers ./internal/service

swag init:
	swag init -g cmd/main.go

