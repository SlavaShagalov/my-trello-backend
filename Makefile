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

#.PHONY: stop
#stop:
#	docker compose -f docker-compose.yml stop

.PHONY: up
up:
	make api-up
	make monitoring-up

.PHONY: stop
stop:
	make api-stop
	make monitoring-stop

.PHONY: down
down:
	docker compose down -v
	sudo rm -rf ./postgres/primary/pgdata
	sudo rm -rf ./postgres/primary/archive
	sudo rm -rf ./postgres/standby/pgdata

.PHONY: jaeger
jaeger:
	docker run --name jaeger \
	-e COLLECTOR_OTLP_ENABLED=true \
	-e OTEL_EXPORTER_OTLP_ENDPOINT=http://0.0.0.0:4318/v1/traces \
	-p 16686:16686 \
	-p 4317:4317 \
	-p 4318:4318 \
	jaegertracing/all-in-one:1.35

#	docker run --rm --name jaeger \
#	-p 16686:16686 \
#	-p 4318:4318 \
#	-e OTEL_EXPORTER_OTLP_ENDPOINT="http://127.0.0.1:4318" \
#	jaegertracing/example-hotrod:1.53 \
#	all --otel-exporter=otlp

.PHONY: api-up
api-up:
	docker compose -f docker-compose.yml up -d --build db sessions-db api-main balancer

.PHONY: api-stop
api-stop:
	docker compose -f docker-compose.yml stop db sessions-db api-main balancer

.PHONY: monitoring-up
monitoring-up:
	docker compose -f docker-compose.yml up -d --build node-exporter prometheus grafana jaeger

.PHONY: monitoring-stop
monitoring-stop:
	docker compose -f docker-compose.yml stop node-exporter prometheus grafana jaeger

# ===== LOGS =====

service = node-exporter
.PHONY: logs
logs:
	docker compose logs -f "$(service)"

name = main
.PHONY: logs-api
logs-api:
	tail -f -n +1 "cmd/api/logs/$(name).log" | batcat --paging=never --language=log

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
	go test -shuffle=on ./tests/unit/...

.PHONY: integration-test
integration-test:
	go test ./tests/integration/...
	#go test -count=50 -bench ./tests/integration/...

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
