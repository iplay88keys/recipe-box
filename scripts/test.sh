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
            skipIntegration=true
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
    ./scripts/start_database.sh > /dev/null 2>&1

    function finish {
      ./scripts/stop_database.sh > /dev/null 2>&1
    }
    trap finish EXIT

    ./scripts/migrate_database.sh
fi

if [[ "${skipBackend}" = "false" ]]; then
    echo "Running ginkgo"
    ginkgo -r -p
fi
