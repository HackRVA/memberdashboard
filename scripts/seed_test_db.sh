#!/bin/bash

##
#
# This script assumes you're already running a local postgres db.
# you can stand one up with:
#   docker compose -f ./deployments/docker-compose.yml up -d
#
##

export DB_CONNECTION_STRING=postgresql://test:test@localhost:5432/membership

echo "running `make migrate-up` to build schema"
make migrate-up

echo "running `make seed` to seed db with test data"
make seed

echo "user test@test.com with password test "

