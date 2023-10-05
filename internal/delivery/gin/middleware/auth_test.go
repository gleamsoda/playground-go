package middleware

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/require"

// 	"playground/internal/delivery/gin/helper"
// 	"playground/internal/pkg/token"
// )

// func TestAuthMiddleware(t *testing.T) {
// 	testCases := []struct {
// 		name          string
// 		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Manager)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
// 				helper.AddAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, "user", time.Minute)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "NoAuthorization",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "UnsupportedAuthorization",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
// 				addAuthorization(t, request, tokenMaker, "unsupported", "user", time.Minute)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidAuthorizationFormat",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
// 				addAuthorization(t, request, tokenMaker, "", "user", time.Minute)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "ExpiredToken",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
// 				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, "user", -time.Minute)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			server := newTestServer(t, nil)
// 			authPath := "/auth"
// 			server.server.Handler.(*gin.Engine).GET(
// 				authPath,
// 				Auth(server.tm),
// 				func(ctx *gin.Context) {
// 					ctx.JSON(http.StatusOK, gin.H{})
// 				},
// 			)

// 			recorder := httptest.NewRecorder()
// 			request, err := http.NewRequest(http.MethodGet, authPath, nil)
// 			require.NoError(t, err)

// 			tc.setupAuth(t, request, server.tm)
// 			server.server.Handler.(*gin.Engine).ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)
// 		})
// 	}
// }
