#!/usr/bin/env bash

docker info > /dev/null 2>&1
exit_code=$?
if [[ "${exit_code}" -ne 0 ]]; then
    echo "Could not reach docker daemon Is docker installed and running?"
    exit 1
fi

set -e

echo "Stopping the database if currently running"
./scripts/stop_database.sh

pushd db
    echo "Bringing up the new database"
    docker-compose up & > /dev/null 2>&1
popd

exit_code=1
set +e
while [[ "${exit_code}" -eq 1 ]]; do
    echo "Waiting for mysql to be available..."
    mysqladmin -u "${DATABASE_USERNAME}" \
        -p"${DATABASE_PASSWORD}" \
        -h "${DATABASE_HOST}" \
        -P "${DATABASE_PORT}" ping  > /dev/null 2>&1

    exit_code=$?

    sleep 5
done
set -e
