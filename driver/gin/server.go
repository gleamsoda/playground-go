package gin

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"

	"playground/app/repository"
	"playground/app/usecase"
	"playground/config"
	"playground/pkg/token"
)

func NewServer(cfg config.Config) (*http.Server, error) {
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
	tm, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	// repositories
	r := repository.NewRepository(conn)
	// usecases
	u := usecase.NewUsecase(r, tm, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)
	// handlers
	h := NewHandler(u)

	svr := gin.Default()
	svr.GET("/health", health(conn))
	svr.POST("/users", h.createUser)
	svr.POST("/login", h.loginUser)
	svr.POST("/tokens/renew_access", h.renewAccessToken)

	auth := svr.Group("/").Use(authMiddleware(tm))
	auth.POST("/accounts", h.createAccount)
	auth.GET("/accounts/:id", h.getAccount)
	auth.GET("/accounts", h.listAccounts)
	auth.POST("/transfers", h.createTransfer)

	return &http.Server{
		Addr:    cfg.HTTPServerAddress,
		Handler: svr,
	}, nil
}
