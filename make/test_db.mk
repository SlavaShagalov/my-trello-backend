.PHONY: test-db-up
test-db-up:
	docker compose -f docker-compose.yml up -d judi-test-db

.PHONY: test-db-stop
test-db-stop:
	docker compose -f docker-compose.yml stop judi-test-db

.PHONY: test-db-down
test-db-down:
	docker compose -f docker-compose.yml stop judi-test-db
	docker compose -f docker-compose.yml rm -f judi-test-db

.PHONY: test-db-create-schema
test-db-create-schema:
	./scripts/db/create_test_db_schema.sh

.PHONY: test-db-fill
test-db-fill:
	./scripts/db/fill_test_data.sh

.PHONY: test-db-prepare
test-db-prepare: test-db-create-schema test-db-fill
