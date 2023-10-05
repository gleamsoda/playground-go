package helper

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"playground/internal/pkg/token"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func AddAuthorization(
	t *testing.T,
	req *http.Request,
	tm token.Manager,
	authorizationType string,
	username string,
	d time.Duration,
) {
	t.Helper()
	tkn, payload, err := tm.Create(username, d)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, tkn)
	req.Header.Set(AuthorizationHeaderKey, authorizationHeader)
}
