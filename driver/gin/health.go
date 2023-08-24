package gin

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func health(conn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}
