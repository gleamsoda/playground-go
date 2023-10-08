package helper

import (
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
