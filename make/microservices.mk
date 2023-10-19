.PHONY: auth-up
auth-up:
	docker compose -f docker-compose.yml up -d --build judi-auth

.PHONY: auth-stop
auth-stop:
	docker compose -f docker-compose.yml stop judi-auth

.PHONY: auth-down
auth-down:
	docker compose -f docker-compose.yml stop judi-auth
	docker compose -f docker-compose.yml rm -f judi-auth

.PHONY: workspaces-up
workspaces-up:
	docker compose -f docker-compose.yml up -d --build judi-workspaces

.PHONY: workspaces-stop
workspaces-stop:
	docker compose -f docker-compose.yml stop judi-workspaces

.PHONY: workspaces-down
workspaces-down:
	docker compose -f docker-compose.yml stop judi-workspaces
	docker compose -f docker-compose.yml rm -f judi-workspaces
