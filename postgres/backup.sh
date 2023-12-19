#!/bin/bash

sleep 2;
pg_basebackup -h db -p 5432 -U replicator -D /var/lib/postgresql/data -Fp -Xs -R;
