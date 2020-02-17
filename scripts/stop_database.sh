#!/usr/bin/env bash

set -e

pushd db
  docker-compose down
popd
