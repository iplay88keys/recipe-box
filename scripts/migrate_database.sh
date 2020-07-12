#!/usr/bin/env bash

set -e

pushd migrations
    echo "Migrating the database"
    flyway migrate \
        -url="jdbc:mysql://${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}" \
        -user="${DATABASE_USERNAME}" \
        -password="${DATABASE_PASSWORD}" \
        -locations=filesystem:.

    echo "Importing example data into the database"
    mysql -u "${DATABASE_USERNAME}" \
        -p"${DATABASE_PASSWORD}" \
        -h "${DATABASE_HOST}" \
        -P "${DATABASE_PORT}" \
        -D "${DATABASE_NAME}" < examples/example.sql
popd
