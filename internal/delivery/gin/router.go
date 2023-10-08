package gin

import (
	"github.com/gin-gonic/gin"

	"playground/internal/delivery/gin/handler"
	"playground/internal/delivery/gin/middleware"
)

func NewRouter(h *handler.Handler, middlewareAuth gin.HandlerFunc) *gin.Engine {
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

	auth := router.Group("/").Use(middlewareAuth)
	auth.POST("/accounts", h.CreateAccount)
	auth.GET("/accounts/:id", h.GetAccount)
	auth.GET("/accounts", h.ListAccounts)
	auth.POST("/transfers", h.CreateTransfer)
	return router
}
