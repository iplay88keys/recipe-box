#!/usr/bin/env bash

set -e

. ./scripts/dev_db_creds.sh

pushd migrations
    echo "Importing example data into the database"
    mysql -u "${DATABASE_USERNAME}" \
        -p"${DATABASE_PASSWORD}" \
        -h "${DATABASE_HOST}" \
        -P "${DATABASE_PORT}" \
        -D "${DATABASE_NAME}" < examples/example.sql
popd
