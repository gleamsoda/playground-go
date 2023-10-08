package gqlgen

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"playground/internal/delivery/gqlgen/helper"
	"playground/internal/pkg/token"
)

func MetadataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md := &helper.Metadata{
			UserAgent: r.UserAgent(),
			ClientIP:  r.RemoteAddr,
		}
		ctx := context.WithValue(r.Context(), helper.MetadataCtxkey, md)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthMiddlewareFunc(tm token.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get(helper.AuthorizationHeaderKey)
			if len(authorizationHeader) == 0 {
				next.ServeHTTP(w, r)
				return
			}
			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				next.ServeHTTP(w, r)
				return
			}
			authorizationType := strings.ToLower(fields[0])
			if authorizationType != helper.AuthorizationTypeBearer {
				next.ServeHTTP(w, r)
				return
			}
			accessToken := fields[1]
			payload, err := tm.Verify(accessToken)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, helper.ErrorResponse(err))
				return
			}
			ctx := context.WithValue(r.Context(), helper.AuthCtxkey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
