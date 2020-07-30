include .env
include .test.env
export

.PHONY: openapi
openapi: openapi_http openapi_js

.PHONY: openapi_http
openapi_http:
	oapi-codegen -generate types -o internal/trainings/openapi_types.gen.go -package main api/openapi/trainings.yml
	oapi-codegen -generate chi-server -o internal/trainings/openapi_api.gen.go -package main api/openapi/trainings.yml

	oapi-codegen -generate types -o internal/trainer/openapi_types.gen.go -package main api/openapi/trainer.yml
	oapi-codegen -generate chi-server -o internal/trainer/openapi_api.gen.go -package main api/openapi/trainer.yml

	oapi-codegen -generate types -o internal/users/openapi_types.gen.go -package main api/openapi/users.yml
	oapi-codegen -generate chi-server -o internal/users/openapi_api.gen.go -package main api/openapi/users.yml

.PHONY: openapi_js
openapi_js:
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.0 generate \
        -i /local/api/openapi/trainings.yml \
        -g javascript \
        -o /local/web/src/repositories/clients/trainings

	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.0 generate \
		-i /local/api/openapi/trainer.yml \
		-g javascript \
		-o /local/web/src/repositories/clients/trainer

	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.0 generate \
		-i /local/api/openapi/users.yml \
		-g javascript \
		-o /local/web/src/repositories/clients/users

.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:internal/common/genproto/trainer -I api/protobuf api/protobuf/trainer.proto
	protoc --go_out=plugins=grpc:internal/common/genproto/users -I api/protobuf api/protobuf/users.proto

.PHONY: lint
lint:
	@./scripts/lint.sh trainer
	@./scripts/lint.sh trainings
	@./scripts/lint.sh users

.PHONY: mycli
mycli:
	mycli -u ${MYSQL_USER} -p ${MYSQL_PASSWORD} ${MYSQL_DATABASE}

INERNAL_PACKAGES := $(wildcard internal/*)

ifeq (test,$(firstword $(MAKECMDGOALS)))
  TEST_ARGS := $(subst $$,$$$$,$(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS)))
  $(eval $(TEST_ARGS):;@:)
endif

.PHONY: test $(INERNAL_PACKAGES)
test: $(INERNAL_PACKAGES)
$(INERNAL_PACKAGES):
	@(cd $@ && go test -count=1 -race ./... $(subst $$$$,$$,$(TEST_ARGS)))
