include .env
export

.PHONY: openapi
openapi: openapi_http openapi_js

.PHONY: openapi_http
openapi_http:
	oapi-codegen -generate types -o internal/trainings/ports/openapi_types.gen.go -package ports api/openapi/trainings.yml
	oapi-codegen -generate chi-server -o internal/trainings/ports/openapi_api.gen.go -package ports api/openapi/trainings.yml
	oapi-codegen -generate types -o internal/common/client/trainings/openapi_types.gen.go -package trainings api/openapi/trainings.yml
	oapi-codegen -generate client -o internal/common/client/trainings/openapi_client_gen.go -package trainings api/openapi/trainings.yml

	oapi-codegen -generate types -o internal/trainer/ports/openapi_types.gen.go -package ports api/openapi/trainer.yml
	oapi-codegen -generate chi-server -o internal/trainer/ports/openapi_api.gen.go -package ports api/openapi/trainer.yml
	oapi-codegen -generate types -o internal/common/client/trainer/openapi_types.gen.go -package trainer api/openapi/trainer.yml
	oapi-codegen -generate client -o internal/common/client/trainer/openapi_client_gen.go -package trainer api/openapi/trainer.yml

	oapi-codegen -generate types -o internal/users/openapi_types.gen.go -package main api/openapi/users.yml
	oapi-codegen -generate chi-server -o internal/users/openapi_api.gen.go -package main api/openapi/users.yml
	oapi-codegen -generate types -o internal/common/client/users/openapi_types.gen.go -package users api/openapi/users.yml
	oapi-codegen -generate client -o internal/common/client/users/openapi_client_gen.go -package users api/openapi/users.yml

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
	protoc \
		--proto_path=api/protobuf api/protobuf/trainer.proto \
		--go_out=internal/common/genproto/trainer --go_opt=paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
		--go-grpc_out=internal/common/genproto/trainer --go-grpc_opt=paths=source_relative \

	protoc \
		--proto_path=api/protobuf api/protobuf/users.proto \
		--go_out=internal/common/genproto/users --go_opt=paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
		--go-grpc_out=internal/common/genproto/users --go-grpc_opt=paths=source_relative \

.PHONY: lint
lint:
	@./scripts/lint.sh trainer
	@./scripts/lint.sh trainings
	@./scripts/lint.sh users

.PHONY: fmt
fmt:
	goimports -l -w internal/

.PHONY: mycli
mycli:
	mycli -u ${MYSQL_USER} -p ${MYSQL_PASSWORD} ${MYSQL_DATABASE}

test:
	@./scripts/test.sh common .e2e.env
	@./scripts/test.sh trainer .test.env
	@./scripts/test.sh trainings .test.env
	@./scripts/test.sh users .test.env
