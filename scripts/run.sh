#!/usr/bin/env bash

set -e

buildUI=true
while test $# -gt 0
do
    case "$1" in
        --skip-ui-build)
            buildUI=false
            ;;
        --*)
            echo "bad option $1"
            exit 0
            ;;
        *)
            echo "bad argument $1"
            exit 0
            ;;
    esac
    shift
done

if [[ "${buildUI}" = "true" ]]; then
    pushd ui
        yarn build
    popd
fi

export DATABASE_USERNAME="user"
export DATABASE_PASSWORD="password"
export DATABASE_HOST="127.0.0.1"
export DATABASE_PORT="3306"
export DATABASE_NAME="db"

#echo "Restarting the database"
#./scripts/start_database.sh
#
#echo "Migrating the database"
#./scripts/migrate_database.sh

#echo "Importing example data"
#./scripts/import_example_database_data.sh
#
#function finish {
#    echo "Stopping the database"
#    ./scripts/stop_database.sh
#}
#trap finish EXIT

echo "Exporting env vars"
export DATABASE_URL="mysql://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@tcp(${DATABASE_HOST}:${DATABASE_PORT})/${DATABASE_NAME}"
export REDIS_URL="127.0.0.1:6379"
export ACCESS_SECRET="access_secret"
export REFRESH_SECRET="refresh_secret"

echo "Starting the app"
go run main.go
