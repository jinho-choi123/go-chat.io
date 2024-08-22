
dev:
	air

db-up:
	docker compose -p go_chat_io -f .docker/docker-compose.yml --env-file .config/.env.development up

db-down:
	docker compose -p go_chat_io -f .docker/docker-compose.yml --env-file .config/.env.development down -v