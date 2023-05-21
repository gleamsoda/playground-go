package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gleamsoda/go-playground/cmd/gin/internal"
	"github.com/gleamsoda/go-playground/config"
	"github.com/gleamsoda/go-playground/repo"
	"github.com/gleamsoda/go-playground/usecase"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := sql.Open("mysql", cfg.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	entryRepo := repo.NewEntryRepository(conn)
	entryUsecase := usecase.NewEntryUsecase(entryRepo)
	entryHandler := internal.NewEntryHandler(entryRepo, entryUsecase)

	r.GET("/ping", ping(conn))
	r.POST("/entries", entryHandler.Create)
	r.GET("/entries/:id", entryHandler.Get)
	log.Fatal(r.Run())
}

func ping(conn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := conn.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
