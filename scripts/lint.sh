#!/bin/bash
set -e

readonly service="$1"

cd "./pkg/$service"
go vet .
