package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"playground/internal/delivery/gin/helper"
	"playground/internal/pkg/token"
	"playground/internal/wallet"
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
	if a, err := h.w.CreateAccount(ctx, &wallet.CreateAccountParams{
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
	if a, err := h.w.GetAccount(ctx, &wallet.GetAccountsParams{
		ID:    req.ID,
		Owner: authPayload.Username,
	}); err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, a)
	}
}

type listAccountsRequest struct {
	Limit  int32 `form:"limit" binding:"required,min=1,max=100"`
	Offset int32 `form:"offset"`
}

func (h *Handler) ListAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(helper.AuthorizationPayloadKey).(*token.Payload)
	if as, err := h.w.ListAccounts(ctx, &wallet.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.Limit,
		Offset: req.Offset,
	}); err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, as)
	}
}
