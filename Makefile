include make/db.mk
include make/test_db.mk
include make/redis.mk
include make/test_redis.mk

EASYJSON_PATHS = ./internal/...

# ===== RUN =====

.PHONY: start
start:
	make format
	make swag
	#docker compose -f docker-compose.yml up -d --build data-storage-rep ds-admin api-main api-read-1 api-read-2 api-mirror balancer
	docker compose -f docker-compose.yml up -d --build db db-repl ds-admin api-main api-read-1 api-read-2 api-mirror balancer
	#docker compose -f docker-compose.yml up -d --build ds-admin api-main balancer

.PHONY: stop
stop:
	docker compose -f docker-compose.yml stop

# ===== LOGS =====

name = main
.PHONY: api-logs
api-logs:
	tail -f -n +1 "cmd/api/logs/$(name).log" | batcat --paging=never --language=log

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

.PHONY: run-test-containers
run-test-containers:
	#docker compose -f docker-compose.yml up -d --build db sessions-storage api-main balancer test
	docker compose -f docker-compose.yml up -d --build test-db test-sessions-db test-api test

.PHONY: unit-test
unit-test:
	ALLURE_OUTPUT_PATH=$(CURDIR) go test ./tests/unit/...

.PHONY: integration-test
integration-test:
	go test ./tests/integration/...
	#go test -count=50 ./tests/integration/...

.PHONY: e2e-test
e2e-test:
	go test ./tests/e2e/...

.PHONY: unit-cover
unit-cover:
	go test -covermode=atomic -coverprofile=cover.out ./internal/...
	go tool cover -func=cover.out
	go tool cover -html=cover.out -o coverage.html
	@rm cover.out

.PHONY: integration-cover
integration-cover:
	./scripts/db/run_integration_cover.sh

suite = auth
.PHONY: test-logs
test-logs:
	tail -f -n +1 "tests/logs/$(suite).log" | batcat --paging=never --language=log
