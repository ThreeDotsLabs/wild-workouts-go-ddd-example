FROM golang:1.16.0 AS go
# TODO: add openapitools/openapi-generator-cli

FROM go AS protoc
# TODO: update to 3.15.3
RUN VERSION=3.13.0 && \
    apt-get update && \
    apt-get install unzip && \
    wget -O x.zip https://github.com/protocolbuffers/protobuf/releases/download/v$VERSION/protoc-$VERSION-linux-x86_64.zip && \
    unzip x.zip -d protoc

FROM go AS protobuf
# TODO: update to v1.4.3
RUN go install github.com/golang/protobuf/protoc-gen-go@v1.3.3

# TODO: migrate from github.com/golang/protobuf => google.golang.org/protobuf + google.golang.org/grpc

FROM go AS oapi-codegen
# TODO: update to v1.5.1
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.3.9

FROM go
COPY --from=protoc       /go/protoc/                /usr/
COPY --from=protobuf     /go/bin/protoc-gen-go      /usr/bin/
COPY --from=oapi-codegen /go/bin/oapi-codegen       /usr/bin/
