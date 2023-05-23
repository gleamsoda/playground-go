package internal

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(conn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}
