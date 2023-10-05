package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"playground/internal/delivery/gin/helper"
	"playground/internal/pkg/token"
)

// Auth creates a gin middleware for authorization
func Auth(tm token.Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(helper.AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != helper.AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tm.Verify(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.ErrorResponse(err))
			return
		}

		ctx.Set(helper.AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
