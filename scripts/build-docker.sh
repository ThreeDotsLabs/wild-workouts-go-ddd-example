#!/bin/bash
readonly service="$1"
readonly project_id="$2"

docker build -t "gcr.io/$project_id/$service" "./internal" -f "./docker/app-prod/Dockerfile" --build-arg "SERVICE=$service"
docker push "gcr.io/$project_id/$service"
