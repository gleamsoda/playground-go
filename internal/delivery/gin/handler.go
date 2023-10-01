package gin

import (
	"github.com/gin-gonic/gin"

	"playground/app"
)

type handler struct {
	u app.Usecase
}

func NewHandler(u app.Usecase, authMiddleware gin.HandlerFunc) *gin.Engine {
	r := &handler{u: u}

	svr := gin.Default()
	svr.GET("/health", health())
	svr.POST("/users", r.createUser)
	svr.POST("/users/login", r.loginUser)
	svr.POST("/tokens/renew_access", r.renewAccessToken)

	auth := svr.Group("/").Use(authMiddleware)
	auth.POST("/accounts", r.createAccount)
	auth.GET("/accounts/:id", r.getAccount)
	auth.GET("/accounts", r.listAccounts)
	auth.POST("/transfers", r.createTransfer)
	return svr
}
