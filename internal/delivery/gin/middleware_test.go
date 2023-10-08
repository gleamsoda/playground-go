package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"playground/internal/config"
	"playground/internal/delivery/gin/middleware"
	"playground/internal/pkg/token"
)

func TestAuthMiddleware(t *testing.T) {
	tm, err := token.NewPasetoManager(config.Get().TokenSymmetricKey)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name          string
		setup         func(t *testing.T, request *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setup: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, tm, "bearer", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setup: func(t *testing.T, request *http.Request) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setup: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, tm, "unsupported", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setup: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, tm, "", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setup: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, tm, "bearer", "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/", middleware.Auth(tm), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			tc.setup(t, request)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
