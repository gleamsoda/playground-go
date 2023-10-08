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

	"playground/internal/config"
	"playground/internal/delivery/gin/handler"
	"playground/internal/delivery/gin/helper"
	"playground/internal/pkg/mail"
	"playground/internal/pkg/token"
	"playground/internal/wallet"
	"playground/internal/wallet/usecase"
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
	do.Provide(injector, handler.NewHandler)
	do.Provide(injector, usecase.NewUsecase)
	do.ProvideValue[wallet.Repository](injector, nil)
	do.ProvideValue[wallet.Dispatcher](injector, nil)
	do.ProvideValue[mail.Sender](injector, nil)
	do.ProvideValue[token.Manager](injector, tm)
	do.ProvideNamedValue(injector, "AccessTokenDuration", cfg.AccessTokenDuration)
	do.ProvideNamedValue(injector, "RefreshTokenDuration", cfg.RefreshTokenDuration)
	return injector
})

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
