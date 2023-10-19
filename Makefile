include make/db.mk
include make/test_db.mk
include make/redis.mk
include make/test_redis.mk
include make/microservices.mk

EASYJSON_PATHS = ./internal/...

.PHONY: deploy
deploy:
	docker compose -f docker-compose.yml up -d --build api

.PHONY: stop
stop:
	docker compose -f docker-compose.yml stop

.PHONY: api-logs
api-logs:
	tail -f cmd/api/logs/api.log | batcat --paging=never --language=log

.PHONY: mocks
mocks:
	./scripts/gen_mocks.sh

.PHONY: easyjson
easyjson:
	go generate ${EASYJSON_PATHS}

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
