#!/bin/bash

# DB up
make test-db-up
make test-redis-up
sleep 3
make test-db-prepare

# RUN tests
# -count=100 - check unstable tests
go test -count=1 -coverpkg=./.../usecase/...,./.../repository/... -covermode=atomic \
  -coverprofile=cover.out ./tests/integration/... &&
  go tool cover -func=cover.out
go tool cover -html=cover.out -o coverage.html
rm cover.out

# DB down
trap 'make test-db-down; make test-redis-down' EXIT
