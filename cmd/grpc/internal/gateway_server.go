package internal

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"playground/cmd/grpc/internal/boundary"
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
	if err := boundary.RegisterPlaygroundHandlerFromEndpoint(ctx, mux, grpcServerAddress, opts); err != nil {
		return nil, fmt.Errorf("cannot register handler server: %w", err)
	}

	return mux, nil
}
