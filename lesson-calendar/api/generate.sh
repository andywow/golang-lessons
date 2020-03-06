#!/bin/bash

set -e

BASEDIR="$(dirname $0)"

echo "Generating go code for api proto"
protoc ${BASEDIR}/api.proto -I${BASEDIR} --go_out=plugins=grpc,paths=source_relative:${BASEDIR}/../pkg/eventapi
echo "Done"
