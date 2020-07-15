#!/usr/bin/env bash

set -e

. ./scripts/dev_db_creds.sh

skipUI=false
skipIntegration=false
skipBackend=false

while test $# -gt 0; do
    case "$1" in
        --skip-ui)
            skipUI=true
            ;;
        --skip-backend)
            skipBackend=true
            ;;
        --skip-integration)
            skipIntegration=true
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

if [[ "${skipUI}" = "false" ]]; then
    pushd ui
        echo "Running 'yarn test'"
        yarn test --watchAll=false

        echo "Compiling the UI"
        yarn build
    popd
fi

if [[ "${skipIntegration}" = "false" ]]; then
    exit_code=1
    set +e
    echo "Checking to see if mysql is available"
    mysqladmin -u "${DATABASE_USERNAME}" \
        -p"${DATABASE_PASSWORD}" \
        -h "${DATABASE_HOST}" \
        -P "${DATABASE_PORT}" ping  > /dev/null 2>&1

    exit_code=$?
    set -e

    if [[ "${exit_code}" -eq 1 ]]; then
        echo "mysql is not running, starting it for testing"
        ./scripts/start_database.sh > /dev/null 2>&1

        function finish {
          ./scripts/stop_database.sh > /dev/null 2>&1
        }
        trap finish EXIT
    else
        ./scripts/clean_database.sh
    fi

    ./scripts/migrate_database.sh
fi


echo "Setting required env vars"
export REDIS_URL="redis://:@127.0.0.1:6379"
export ACCESS_SECRET="access_secret"
export REFRESH_SECRET="refresh_secret"

if [[ "${skipBackend}" = "false" ]]; then
    echo "Running ginkgo for everything except integration"
    ginkgo -r -p -skipPackage pkg/integration
fi

if [[ "${skipIntegration}" = "false" ]]; then
    echo "Running ginkgo for integration"
    ginkgo -r pkg/integration
fi
