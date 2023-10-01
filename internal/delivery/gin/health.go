package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}
