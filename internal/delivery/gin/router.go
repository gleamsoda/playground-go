package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"

	"playground/internal/delivery/gin/handler"
	"playground/internal/delivery/gin/middleware"
	"playground/internal/pkg/token"
)

func NewRouter(i *do.Injector) (*gin.Engine, error) {
	h := do.MustInvoke[*handler.Handler](i)
	tm := do.MustInvoke[token.Manager](i)
	router := gin.New()
	router.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.ErrorHandler(),
	)

	router.GET("/health", handler.Health())
	router.POST("/users", h.CreateUser)
	router.POST("/users/login", h.LoginUser)
	router.POST("/tokens/renew_access", h.RenewAccessToken)

	auth := router.Group("/").Use(middleware.Auth(tm))
	auth.POST("/accounts", h.CreateAccount)
	auth.GET("/accounts/:id", h.GetAccount)
	auth.GET("/accounts", h.ListAccounts)
	auth.POST("/transfers", h.CreateTransfer)
	return router, nil
}
