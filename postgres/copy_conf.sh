#!/bin/bash

cp /configs/postgresql.conf /var/lib/postgresql/data/postgresql.conf
mkdir -p /var/lib/postgresql/data/archive
cp /configs/pg_hba.conf /var/lib/postgresql/data/pg_hba.conf
