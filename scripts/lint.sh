#!/bin/bash
set -e

if [ "$2" == "-install" ]; then
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.43.0
  go install github.com/roblaszczak/go-cleanarch@latest
fi

readonly service="$1"

cd "./internal/$service"
golangci-lint run