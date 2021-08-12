#!/bin/bash

export FIRESTORE_EMULATOR_HOST=fake:8787
export TRAINER_GRPC_ADDR=fake:3000
export USERS_GRPC_ADDR=fake:3001

rm -rf out
mkdir out

go run .

if ! type "plantuml" > /dev/null; then
  echo "Please install plantuml to generate PNG diagrams automatically"
fi

for f in out/*.plantuml; do plantuml "$f"; done
