//go:build tools

package tools

import (
	// imports for grpc-ecosystem/grpc-gateway
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"

	// imports for golang/mock
	_ "go.uber.org/mock/mockgen/model" // https://github.com/golang/mock#debugging-errors

	// imports for gqlgen
	_ "github.com/99designs/gqlgen"
)
