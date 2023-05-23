package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gleamsoda/go-playground/domain"
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
	var args domain.CreateUserParams
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
