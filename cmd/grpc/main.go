package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gleamsoda/go-playground/cmd/grpc/internal"
	"github.com/gleamsoda/go-playground/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		gw, err := internal.NewGatewayServer(cfg.GRPCServerAddress)
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
		s, err := internal.NewServer(cfg)
		if err != nil {
			log.Fatal(err)
		}
		l, err := net.Listen("tcp", cfg.GRPCServerAddress)
		if err != nil {
			log.Fatal("cannot create listener:", err)
		}
		log.Printf("start gRPC server at %s", cfg.GRPCServerAddress)
		log.Fatal(s.Serve(l))
	}()
}
