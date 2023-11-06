.PHONY: db-start
db-start:
	mkdir -p ./postgres/primary/archive
	sudo chown 999:999 ./postgres/primary/archive
	docker compose -f docker-compose.yml up -d --build db
	docker run -it --rm \
    --net trello_network \
    -v ./postgres/standby/pgdata:/var/lib/postgresql/data \
    -v ./postgres/backup.sh:/backup.sh \
    --entrypoint /bin/bash postgres /backup.sh
	docker compose -f docker-compose.yml up -d --build db-repl

.PHONY: db-down
db-down:
	docker compose down -v
	sudo rm -rf ./postgres/primary/pgdata
	sudo rm -rf ./postgres/primary/archive
	sudo rm -rf ./postgres/standby/pgdata
