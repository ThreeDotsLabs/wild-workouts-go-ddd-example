#!/bin/bash
set -e

readonly service="$1"

cd "./internal/$service"
go vet .
