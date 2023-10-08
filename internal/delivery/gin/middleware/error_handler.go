package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/morikuni/failure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"playground/internal/pkg/apperr"
)

// ErrorHandler creates a gin middleware for error handling
func ErrorHandler() gin.HandlerFunc {
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	return func(ctx *gin.Context) {
		ctx.Next()
		if err := ctx.Errors.Last(); err != nil {
			if code, ok := failure.CodeOf(err.Err); ok {
				cs, _ := failure.CallStackOf(err.Err)
				log.Error().
					Str("line", fmt.Sprintf("%s", cs.HeadFrame())).
					Msg(err.Err.Error())
				switch code {
				case apperr.Internal:
					ctx.JSON(http.StatusInternalServerError, errorResponse(err.Err))
				case apperr.InvalidArgument:
					ctx.JSON(http.StatusBadRequest, errorResponse(err.Err))
				case apperr.NotFound:
					ctx.JSON(http.StatusNotFound, errorResponse(err.Err))
				case apperr.AlreadyExists:
					ctx.JSON(http.StatusConflict, errorResponse(err.Err))
				case apperr.Unauthenticated:
					ctx.JSON(http.StatusUnauthorized, errorResponse(err.Err))
				case apperr.PermissionDenied:
					ctx.JSON(http.StatusForbidden, errorResponse(err.Err))
				default:
					ctx.JSON(http.StatusInternalServerError, errorResponse(err.Err))
				}
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err.Err))
		}
	}
}

func errorResponse(err error) gin.H {
	msg := "something went wrong"
	if m, ok := failure.MessageOf(err); ok {
		msg = m
	}
	return gin.H{"error": msg}
}
