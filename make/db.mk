.PHONY: db-up
db-up:
	docker compose -f docker-compose.yml up -d judi-db

.PHONY: db-stop
db-stop:
	docker compose -f docker-compose.yml stop judi-db

.PHONY: db-down
db-down:
	docker compose -f docker-compose.yml stop judi-db
	docker compose -f docker-compose.yml rm -f judi-db

.PHONY: db-create-schema
db-create-schema:
	./scripts/db/create_schema.sh

.PHONY: db-fill
db-fill:
	./scripts/db/fill.sh

.PHONY: db-prepare
db-prepare: db-create-schema db-fill
