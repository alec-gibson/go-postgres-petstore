#!/bin/bash

# Generate spec file with hash of the openapi spec
oapi-codegen -generate "spec" ./openapi/openapi.yaml > openapi/spec.go

# Generate server files
oapi-codegen -generate "types" -package "api" ./openapi/openapi.yaml > pkg/interface/api/types.go
oapi-codegen -generate "server" -package "api" ./openapi/openapi.yaml > pkg/interface/api/controller.go

# Generate client files
oapi-codegen -generate "types" -package "client" ./openapi/openapi.yaml > client/types.go
oapi-codegen -generate "client" -package "client" ./openapi/openapi.yaml > client/client.go

# Generate code for DB interaction
sqlc generate
