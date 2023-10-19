.PHONY: test-redis-up
test-redis-up:
	docker compose -f docker-compose.yml up -d judi-test-redis

.PHONY: test-redis-stop
test-redis-stop:
	docker compose -f docker-compose.yml stop judi-test-redis

.PHONY: test-redis-down
test-redis-down:
	docker compose -f docker-compose.yml stop judi-test-redis
	docker compose -f docker-compose.yml rm -f judi-test-redis
