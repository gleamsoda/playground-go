package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"playground/cmd/grpc/server"
	"playground/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		gw, err := server.NewGatewayServer(ctx, cfg.GRPCServerAddress)
		if err != nil {
			log.Fatal(err)
		}

		l, err := net.Listen("tcp", cfg.HTTPServerAddress)
		if err != nil {
			log.Fatal("cannot create listener:", err)
		}
		log.Printf("start HTTP gateway server at %s", cfg.HTTPServerAddress)
		log.Fatal(http.Serve(l, gw))
	}()

	func() {
		srv, err := server.NewServer(cfg)
		if err != nil {
			log.Fatal(err)
		}
		l, err := net.Listen("tcp", cfg.GRPCServerAddress)
		if err != nil {
			log.Fatal("cannot create listener:", err)
		}
		log.Printf("start gRPC server at %s", cfg.GRPCServerAddress)
		log.Fatal(srv.Serve(l))
	}()
}
