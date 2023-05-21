package internal

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gleamsoda/go-playground/domain"
)

type entryHandler struct {
	r domain.EntryRepository
	u domain.EntryUsecase
}

func NewEntryHandler(r domain.EntryRepository, u domain.EntryUsecase) entryHandler {
	return entryHandler{
		r: r,
		u: u,
	}
}

func (h entryHandler) Create(c *gin.Context) {
	var args domain.CreateEntryParams
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

func (h entryHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if e, err := h.u.Get(c, id); err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, e)
	}
}
