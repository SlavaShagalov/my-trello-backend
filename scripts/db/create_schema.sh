#!/bin/bash

network=$(docker network ls | grep 'trello_network' | awk '{print $2}')

if [[ -z $network ]];
then
    echo No database network found
    exit 1
else
docker run --rm \
    -e PGPASSWORD=1234\
    --network=$network \
    -v $(pwd)/scripts/migrations/:/scripts/migrations/ \
    postgres \
    psql -h data-storage -U slava -d trello_db -f ./scripts/migrations/schema.sql
fi
