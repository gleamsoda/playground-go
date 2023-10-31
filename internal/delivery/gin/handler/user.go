package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"playground/internal/app"
	"playground/internal/delivery/gin/helper"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse(err))
		return
	}

	if e, err := h.w.CreateUser(ctx, &app.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Email:    req.Email,
	}); err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, e)
	}
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *Handler) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse(err))
		return
	}

	if param, err := h.w.LoginUser(ctx, &app.LoginUserParams{
		Username:  req.Username,
		Password:  req.Password,
		UserAgent: ctx.Request.UserAgent(),
		ClientIP:  ctx.ClientIP(),
	}); err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, param)
	}
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) RenewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse(err))
		return
	}

	if resp, err := h.w.RenewAccessToken(ctx, req.RefreshToken); err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, resp)
	}
}
