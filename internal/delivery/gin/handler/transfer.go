package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"playground/internal/app"
	"playground/internal/delivery/gin/helper"
	"playground/internal/pkg/token"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (h *Handler) CreateTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(helper.AuthorizationPayloadKey).(*token.Payload)
	if e, err := h.w.CreateTransfer(ctx, &app.CreateTransferParams{
		RequestUsername: authPayload.Username,
		FromAccountID:   req.FromAccountID,
		ToAccountID:     req.ToAccountID,
		Amount:          req.Amount,
		Currency:        req.Currency,
	}); err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, e)
	}
}
