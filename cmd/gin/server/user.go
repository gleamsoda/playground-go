package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"playground/domain"
)

type userHandler struct {
	u domain.UserUsecase
}

func NewUserHandler(u domain.UserUsecase) userHandler {
	return userHandler{
		u: u,
	}
}

func (h userHandler) Create(c *gin.Context) {
	var args domain.CreateUserInputParams
	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if e, err := h.u.Create(c, args); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, e)
	}
}

func (h userHandler) Get(c *gin.Context) {
	if e, err := h.u.GetByUsername(c, c.Param("username")); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, e)
	}
}

func (h userHandler) Login(ctx *gin.Context) {
	var args domain.LoginUserInputParams
	if err := ctx.ShouldBindJSON(&args); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args.UserAgent = ctx.Request.UserAgent()
	args.ClientIP = ctx.ClientIP()

	if param, err := h.u.Login(ctx, args); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		ctx.JSON(http.StatusOK, param)
	}
}

func (h userHandler) RenewAccessToken(ctx *gin.Context) {
	var args struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := ctx.ShouldBindJSON(&args); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.u.RenewAccessToken(ctx, args.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
