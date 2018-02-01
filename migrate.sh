#!/bin/sh

migrate \
    -path ./migrations \
    -database "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" \
    "$@"