package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gleamsoda/go-playground/cmd/gin/internal"
	"github.com/gleamsoda/go-playground/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	s, err := internal.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Run())
}
