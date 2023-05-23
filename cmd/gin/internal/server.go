package internal

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gleamsoda/go-playground/config"
	"github.com/gleamsoda/go-playground/internal/token"
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

	tm, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	// Repositories
	entryRepo := repo.NewEntryRepository(conn)
	transferRepo := repo.NewTransferRepository(conn)
	walletRepo := repo.NewWalletRepository(conn)
	userRepo := repo.NewUserRepository(conn)
	sessionRepo := repo.NewSessionRepository(conn)

	// usecases
	entryUsecase := usecase.NewEntryUsecase(entryRepo, walletRepo)
	transferUsecase := usecase.NewTransferUsecase(transferRepo, entryRepo, walletRepo, repo.NewTransactionManager(conn))
	walletUsecase := usecase.NewWalletUsecase(walletRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo, tm, cfg)

	// handlers
	entryHandler := NewEntryHandler(entryUsecase)
	transferHandler := NewTransferHandler(transferUsecase)
	walletHandler := NewWalletHandler(walletUsecase)
	userHandler := NewUserHandler(userUsecase)

	s.GET("/health", health(conn))
	s.POST("/users", userHandler.Create)
	s.GET("/users/:username", userHandler.Get)
	s.POST("/login", userHandler.Login)
	s.POST("/tokens/renew_access", userHandler.RenewAccessToken)

	authRoutes := s.Group("/").Use(authMiddleware(tm))
	authRoutes.GET("/wallets", walletHandler.List)
	authRoutes.POST("/wallets", walletHandler.Create)
	authRoutes.GET("/wallets/:id", walletHandler.Get)
	authRoutes.DELETE("/wallets/:id", walletHandler.Delete)
	authRoutes.GET("/wallets/:id/entries", entryHandler.List)
	authRoutes.POST("/transfers", transferHandler.Create)

	return s, nil
}
