package gin

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"playground/internal/app"
	"playground/internal/config"
	"playground/internal/delivery/gin/handler"
	"playground/internal/delivery/gin/helper"
	mock_app "playground/internal/mock/app"
	"playground/internal/pkg/token"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}
	os.Exit(m.Run())
}

// Get reads configuration from file or environment variables.
var GetInjector = sync.OnceValue(func() *do.Injector {
	cfg := config.Get()
	tm, _ := token.NewPasetoManager(cfg.TokenSymmetricKey)

	injector := do.New()
	do.Provide(injector, NewRouter)
	do.Provide(injector, handler.NewHandler)
	do.ProvideValue[app.Repository](injector, nil)
	do.ProvideValue[app.Dispatcher](injector, nil)
	do.ProvideValue[token.Manager](injector, tm)
	do.ProvideNamedValue(injector, "AccessTokenDuration", cfg.AccessTokenDuration)
	do.ProvideNamedValue(injector, "RefreshTokenDuration", cfg.RefreshTokenDuration)
	return injector
})

func NewMockRepository(t *testing.T, ctrl *gomock.Controller) *mock_app.MockRepository {
	mrm := mock_app.NewMockRepository(ctrl)
	mrm.EXPECT().Account().AnyTimes().Return(mock_app.NewMockAccountRepository(ctrl))
	mrm.EXPECT().Transfer().AnyTimes().Return(mock_app.NewMockTransferRepository(ctrl))
	mrm.EXPECT().User().AnyTimes().Return(mock_app.NewMockUserRepository(ctrl))
	mrm.EXPECT().Transaction().AnyTimes().Return(mock_app.NewMockTransaction(ctrl))
	return mrm
}

func addAuthorization(
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
	req.Header.Set(helper.AuthorizationHeaderKey, authorizationHeader)
}
