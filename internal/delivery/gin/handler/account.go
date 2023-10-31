package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"playground/internal/app"
	"playground/internal/delivery/gin/helper"
	"playground/internal/pkg/token"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (h *Handler) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(helper.AuthorizationPayloadKey).(*token.Payload)
	if a, err := h.createAccountUsecase.Execute(ctx, &app.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}); err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, a)
	}
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (h *Handler) GetAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(helper.AuthorizationPayloadKey).(*token.Payload)
	if a, err := h.getAccountUsecase.Execute(ctx, &app.GetAccountsParams{
		ID:    req.ID,
		Owner: authPayload.Username,
	}); err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, a)
	}
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (h *Handler) ListAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(helper.AuthorizationPayloadKey).(*token.Payload)
	if as, err := h.listAccountsUsecase.Execute(ctx, &app.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}); err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, as)
	}
}
