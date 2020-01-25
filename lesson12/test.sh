#!/bin/bash

set -e

TMP_DIR="$(mktemp -d)"

touch "${TMP_DIR}/PATH"
echo TEST24 > "${TMP_DIR}/TEST"

ENV_DIR="$(go run main.go ${TMP_DIR} env)"

rm -rf ${TMP_DIR}

if [[ "$(grep -e '^TEST=TEST24' <<< \"$ENV_DIR\" | wc -l)" != 1 ]]; then
  echo "Error: 'TEST' var no found"
  exit 1
fi

if [[ "$(grep -e '^PATH=' <<< \"$ENV_DIR\" | wc -l)" != 0 ]]; then
  echo "Error: 'PATH' var found"
  exit 1
fi