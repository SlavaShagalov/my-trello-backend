.PHONY: redis-up
redis-up:
	docker compose -f docker-compose.yml up -d judi-redis

.PHONY: redis-stop
redis-stop:
	docker compose -f docker-compose.yml stop judi-redis

.PHONY: redis-down
redis-down:
	docker compose -f docker-compose.yml stop judi-redis
	docker compose -f docker-compose.yml rm -f judi-redis
