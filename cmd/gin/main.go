package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"playground/cmd/gin/server"
	"playground/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.Run())
}
