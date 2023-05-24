package internal

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/gleamsoda/go-playground/cmd/grpc/internal/boundary"
)

// NewServer creates a new gRPC gateway server.
func NewGatewayServer(grpcServerAddress string) (*runtime.ServeMux, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	mux := runtime.NewServeMux(jsonOption)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := boundary.RegisterPlaygroundHandlerFromEndpoint(ctx, mux, grpcServerAddress, opts); err != nil {
		return nil, fmt.Errorf("cannot register handler server: %w", err)
	}

	return mux, nil
}
