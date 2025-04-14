run: down up

up:
	docker compose up -d --build

down:
	docker compose down

tests: unit_test integration_test

unit_test:
	go test -race ./internal/...

integration_test:
	docker compose -f 'test-docker-compose.yaml' up -d
	go test -race ./tests/...
	docker compose -f 'test-docker-compose.yaml' down
