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
	ur := repository.NewUserRepository(conn)
	ar := repository.NewAccountRepository(conn)
	tr := repository.NewTransferRepository(conn)

	// usecases
	uu := usecase.NewUserUsecase(ur, tm, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)
	au := usecase.NewAccountUsecase(ar)
	tu := usecase.NewTransferUsecase(tr, ar)

	// handlers
	uh := NewUserHandler(uu)
	ah := NewAccountHandler(au)
	th := NewTransferHandler(tu)

	svr := gin.Default()
	svr.GET("/health", health(conn))
	svr.POST("/users", uh.createUser)
	svr.POST("/login", uh.login)
	svr.POST("/tokens/renew_access", uh.renewAccessToken)

	auth := svr.Group("/").Use(authMiddleware(tm))
	auth.POST("/accounts", ah.createAccount)
	auth.GET("/accounts/:id", ah.getAccount)
	auth.GET("/accounts", ah.listAccounts)
	auth.POST("/transfers", th.createTransfer)

	return &http.Server{
		Addr:    cfg.HTTPServerAddress,
		Handler: svr,
	}, nil
}
