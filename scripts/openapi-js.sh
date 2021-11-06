#!/bin/bash
set -e

readonly service="$1"

docker run --rm --env "JAVA_OPTS=-Dlog.level=error" -v "${PWD}:/local" \
  "openapitools/openapi-generator-cli:v4.3.0" generate \
  -i "/local/api/openapi/$service.yml" \
  -g javascript \
  -o "/local/web/src/repositories/clients/$service"
