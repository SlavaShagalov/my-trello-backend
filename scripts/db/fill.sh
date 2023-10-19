#!/bin/bash

network=$(docker network ls | grep 'judi_db_network' | awk '{print $2}')

if [[ -z $network ]];
then
    echo No database network found
    exit 1
else 
docker run --rm \
    -e PGPASSWORD=judi_pswd \
    --network=$network \
    -v $(pwd)/scripts/migrations/:/scripts/migrations/ \
    postgres \
    psql -h judi_db -U judi -d judi_db -f ./scripts/migrations/init_data.sql
fi
