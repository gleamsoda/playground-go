package internal

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gleamsoda/go-playground/config"
	"github.com/gleamsoda/go-playground/repo"
	"github.com/gleamsoda/go-playground/usecase"
)

func NewServer(cfg config.Config) (*gin.Engine, error) {
	conn, err := sql.Open("mysql", cfg.DBSource)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	s := gin.Default()

	// Repositories
	entryRepo := repo.NewEntryRepository(conn)
	transferRepo := repo.NewTransferRepository(conn)
	walletRepo := repo.NewWalletRepository(conn)

	// usecases
	entryUsecase := usecase.NewEntryUsecase(entryRepo)
	transferUsecase := usecase.NewTransferUsecase(transferRepo, entryRepo, walletRepo)
	walletUsecase := usecase.NewWalletUsecase(walletRepo)

	// handlers
	entryHandler := NewEntryHandler(entryUsecase)
	transferHandler := NewTransferHandler(transferUsecase)
	walletHandler := NewWalletHandler(walletUsecase)

	s.GET("/health", Health(conn))
	s.POST("/entries", entryHandler.Create)
	s.GET("/entries/:id", entryHandler.Get)
	s.POST("/transfers", transferHandler.Create)
	s.POST("/wallets", walletHandler.Create)
	s.GET("/wallets/:id", walletHandler.Get)
	s.GET("/wallets/:id/entries", entryHandler.List)
	s.GET("/wallets", walletHandler.List)
	s.DELETE("/wallets/:id", walletHandler.Delete)

	return s, nil
}
