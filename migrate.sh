#!/bin/sh

DRIVER=postgres
HOST=localhost
PORT=5432
DATABASE=postgres
USER=postgres
PASSWORD=postgres
SSLMODE=disable

URL="${DRIVER}://${HOST}:${PORT}/${DATABASE}?sslmode=${SSLMODE}&user=${USER}&password=${PASSWORD}"

migrate -path ./database -database ${URL} "$@"
