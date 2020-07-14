#!/usr/bin/env bash

set -e

. ./scripts/dev_db_creds.sh

pushd migrations
    echo "Migrating the database"
    flyway migrate \
        -url="jdbc:mysql://${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}" \
        -user="${DATABASE_USERNAME}" \
        -password="${DATABASE_PASSWORD}" \
        -locations=filesystem:.
popd
