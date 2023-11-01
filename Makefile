include make/db.mk
include make/test_db.mk
include make/redis.mk
include make/test_redis.mk
include make/microservices.mk

EASYJSON_PATHS = ./internal/...

# ===== RUN =====

.PHONY: deploy
deploy:
	make format
	make swag
	docker compose -f docker-compose.yml up -d --build data-storage-rep ds-admin api-main api-read-1 api-read-2 api-mirror balancer

.PHONY: stop
stop:
	docker compose -f docker-compose.yml stop

# ===== LOGS =====

.PHONY: api-main-logs
api-main-logs:
	tail -f -n +1 cmd/api/logs/main.log | batcat --paging=never --language=log

.PHONY: api-1-logs
api-1-logs:
	tail -f -n +1 cmd/api/logs/read_1.log | batcat --paging=never --language=log

.PHONY: api-2-logs
api-2-logs:
	tail -f -n +1 cmd/api/logs/read_2.log | batcat --paging=never --language=log

.PHONY: mirror-logs
mirror-logs:
	tail -f -n +1 cmd/api/logs/mirror.log | batcat --paging=never --language=log

.PHONY: db-logs
db-logs:
	docker compose logs -f data-storage

.PHONY: db-rep-logs
db-rep-logs:
	docker compose logs -f data-storage-rep

.PHONY: balancer-logs
balancer-logs:
	docker compose logs -f balancer

# ===== GENERATORS =====

.PHONY: mocks
mocks:
	./scripts/gen_mocks.sh

.PHONY: easyjson
easyjson:
	go generate ${EASYJSON_PATHS}

.PHONY: swag
swag:
	swag init -g cmd/api/main.go

# ===== FORMAT =====

.PHONY: format
format:
	swag fmt

# ===== TESTS =====

.PHONY: unit-tests
unit-tests:
	go test ./internal/...

.PHONY: integration-tests
integration-tests:
	go test ./tests/integration/...

.PHONY: unit-cover
unit-cover:
	go test -covermode=atomic -coverprofile=cover.out ./internal/...
	go tool cover -func=cover.out
	go tool cover -html=cover.out -o coverage.html
	@rm cover.out

.PHONY: integration-cover
integration-cover:
	./scripts/db/run_integration_cover.sh
