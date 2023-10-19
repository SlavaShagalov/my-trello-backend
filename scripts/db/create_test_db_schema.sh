#!/bin/bash

network=$(docker network ls | grep 'judi_test_db_network' | awk '{print $2}')

if [[ -z $network ]];
then
    echo No database network found
    exit 1
else
docker run --rm \
    -e PGPASSWORD=judi_test_pswd\
    --network=$network \
    -v $(pwd)/scripts/migrations/:/scripts/migrations/ \
    postgres \
    psql -h judi_test_db -U judi_test -d judi_test_db -f ./scripts/migrations/schema.sql
fi
