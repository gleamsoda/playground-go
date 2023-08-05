package server

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"playground/cmd/grpc/server/gen"
)

// NewServer creates a new gRPC gateway server.
func NewGatewayServer(ctx context.Context, grpcServerAddress string) (*runtime.ServeMux, error) {
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	mux := runtime.NewServeMux(jsonOption)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())} // --insecure
	if err := gen.RegisterPlaygroundHandlerFromEndpoint(ctx, mux, grpcServerAddress, opts); err != nil {
		return nil, fmt.Errorf("cannot register handler server: %w", err)
	}

	return mux, nil
}
