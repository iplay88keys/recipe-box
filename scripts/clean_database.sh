#!/usr/bin/env bash

set -e

. ./scripts/dev_db_creds.sh

echo "Clearing Database"
mysql -u "${DATABASE_USERNAME}" \
    -p"${DATABASE_PASSWORD}" \
    -h "${DATABASE_HOST}" \
    -P "${DATABASE_PORT}" \
    -e "DROP DATABASE IF EXISTS ${DATABASE_NAME}; CREATE DATABASE ${DATABASE_NAME};"
