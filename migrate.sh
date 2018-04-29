#!/bin/sh

migrate \
    -path ./migrations \
    -database "postgres://postgres:postgres@artemis:5432/postgres?sslmode=disable" \
    "$@"
