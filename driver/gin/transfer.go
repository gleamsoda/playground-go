package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"playground/app"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (h handler) createTransfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: Get RequestUserID from AuthPayload
	if e, err := h.u.CreateTransfer(c, &app.CreateTransferParams{
		RequestUsername: "example",
		FromAccountID:   req.FromAccountID,
		ToAccountID:     req.ToAccountID,
		Amount:          req.Amount,
		Currency:        req.Currency,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, e)
	}
}
