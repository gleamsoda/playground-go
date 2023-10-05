package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Error(errors.New("test error"))
		c.Status(http.StatusOK)
	}
}
