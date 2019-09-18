.PHONY: up

up:
	@docker-compose -f docker-compose.yml up

build:
	@docker-compose -f docker-compose.yml build