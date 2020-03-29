#!/bin/bash

set -e

if [[ "$DB_HOST" == "" ]]; then
  echo "set DB_HOST variable"
  exit 1
fi

if [[ "$DB_PORT" == "" ]]; then
  echo "set DB_PORT variable"
  exit 1
fi

if [[ "$DB_NAME" == "" ]]; then
  echo "set DB_NAME variable"
  exit 1
fi

if [[ "$DB_USER" == "" ]]; then
  echo "set DB_USER variable"
  exit 1
fi

if [[ "$DB_PASSWORD" == "" ]]; then
  echo "set DB_PASSWORD variable"
  exit 1
fi

BASEDIR="$(readlink -f $(dirname $0))"

docker run --rm --name migrate --network host -v ${BASEDIR}:/migrations migrate/migrate:latest \
  -path=/migrations \
  -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&x-migrations-table=migrate" \
  up
