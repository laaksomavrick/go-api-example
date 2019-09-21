.PHONY: up down build psql run-dev-api run-dev-integration-tests

up:
	@docker-compose -f docker-compose.yml up

down:
	@docker-compose down

build:
	@docker-compose -f docker-compose.yml build

psql:
	@ psql --host=localhost --port=5432 --user=postgres

run-dev-api:
	@GO_ENV=development PORT=3000 POSTGRES_USER=postgres POSTGRES_HOST=0.0.0.0 go run src/main.go

run-dev-integration-tests:
	@API_HOST=0.0.0.0 API_PORT=3000 POSTGRES_USER=postgres POSTGRES_HOST=0.0.0.0 go test -v -count=1 -p=1 tests/messages_test.go


