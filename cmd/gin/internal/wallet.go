package internal

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gleamsoda/go-playground/domain"
)

type walletHandler struct {
	u domain.WalletUsecase
}

func NewWalletHandler(u domain.WalletUsecase) walletHandler {
	return walletHandler{
		u: u,
	}
}

func (h walletHandler) Create(c *gin.Context) {
	var args domain.CreateWalletParams
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

func (h walletHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if e, err := h.u.Get(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, e)
	}
}

func (h walletHandler) List(c *gin.Context) {
	var args domain.ListWalletsParams
	if err := c.ShouldBindQuery(&args); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if es, err := h.u.List(c, args); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, es)
	}
}

func (h walletHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := h.u.Delete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.Status(http.StatusOK)
	}
}
