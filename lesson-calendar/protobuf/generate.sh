#!/bin/bash

set -e

BASEDIR="$(dirname $0)"

echo "Generating go code for event proto"
protoc ${BASEDIR}/event.proto --go_out=plugins=grpc,paths=source_relative:${BASEDIR}/../internal/calendar
echo "Generating go code for api server proto"
protoc ${BASEDIR}/api.proto -I${BASEDIR} --go_out=plugins=grpc,paths=source_relative:${BASEDIR}/../internal/grpcserver
echo "Done"
