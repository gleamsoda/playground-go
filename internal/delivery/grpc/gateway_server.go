package grpc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"playground/internal/config"
	"playground/internal/delivery/grpc/gen"
)

type GatewayServer struct {
	server *http.Server
}

// NewServer creates a new gRPC gateway server.
func NewGatewayServer(ctx context.Context, cfg config.Config) (*GatewayServer, error) {
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
	handler := GatewayLogger(mux)

	return &GatewayServer{
		server: &http.Server{
			Addr:    cfg.HTTPServerAddress,
			Handler: handler,
		},
	}, nil
}

func (s *GatewayServer) Run() error {
	return s.server.ListenAndServe()
}

func (s *GatewayServer) Shutdown() error {
	ctxsd, cancelsd := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelsd()
	if err := s.server.Shutdown(ctxsd); err != nil {
		return fmt.Errorf("gateway server error: failed to shutdown gracefully: %v", err)
	}
	return nil
}

func GatewayLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Msg("received a HTTP request")
	})
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}
