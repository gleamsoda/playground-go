package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gleamsoda/go-playground/domain"
)

type transferHandler struct {
	u domain.TransferUsecase
}

func NewTransferHandler(u domain.TransferUsecase) transferHandler {
	return transferHandler{
		u: u,
	}
}

func (h transferHandler) Create(c *gin.Context) {
	var args domain.CreateTransferParams
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
