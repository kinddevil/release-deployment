NO_COLOR=\033[0m
OK_COLOR=\033[32;01m

SERVICE=campus-backend-service

.PHONY: up down

up:
	@docker-compose up -V --build

down:
	@docker-compose down -v --rmi all && docker images -qf dangling=true | xargs docker rmi
