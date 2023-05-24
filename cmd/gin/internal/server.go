package internal

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gleamsoda/go-playground/config"
	"github.com/gleamsoda/go-playground/internal/token"
	repo "github.com/gleamsoda/go-playground/repo/sqlc"
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
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}
	tm, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	// repositories
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

	svr := gin.Default()
	svr.GET("/health", health(conn))
	svr.POST("/users", userHandler.Create)
	svr.GET("/users/:username", userHandler.Get)
	svr.POST("/login", userHandler.Login)
	svr.POST("/tokens/renew_access", userHandler.RenewAccessToken)

	authRoutes := svr.Group("/").Use(authMiddleware(tm))
	authRoutes.GET("/wallets", walletHandler.List)
	authRoutes.POST("/wallets", walletHandler.Create)
	authRoutes.GET("/wallets/:id", walletHandler.Get)
	authRoutes.DELETE("/wallets/:id", walletHandler.Delete)
	authRoutes.GET("/wallets/:id/entries", entryHandler.List)
	authRoutes.POST("/transfers", transferHandler.Create)

	return svr, nil
}
