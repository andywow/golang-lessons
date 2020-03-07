#!/bin/bash

set -e

if [[ "${GOPATH}" == "" ]]; then
  echo "Specify GOPATH"
  exit 1
fi

go get -u github.com/gogo/protobuf/protoc-gen-gofast

BASEDIR="$(dirname $0)"

protoc ${BASEDIR}/api.proto \
-I=${BASEDIR} \
-I=${GOPATH}/src \
-I=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
--gofast_out=\
Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,\
plugins=grpc,paths=source_relative:${BASEDIR}/../pkg/eventapi
