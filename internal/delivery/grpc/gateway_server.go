package grpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"playground/internal/config"
	"playground/internal/delivery/grpc/gen"
)

// NewServer creates a new gRPC gateway server.
func NewGatewayServer(ctx context.Context, cfg config.Config) (*http.Server, error) {
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
	if err := gen.RegisterPlaygroundHandlerFromEndpoint(ctx, mux, cfg.GRPCServerAddress, opts); err != nil {
		return nil, fmt.Errorf("cannot register handler server: %w", err)
	}
	h := HTTPLogger(mux)

	return &http.Server{
		Addr:    cfg.HTTPServerAddress,
		Handler: h,
	}, nil
}
