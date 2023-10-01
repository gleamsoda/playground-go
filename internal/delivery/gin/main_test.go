package gin

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"

	"playground/internal/config"
	"playground/internal/pkg/token"
	"playground/internal/wallet"
	"playground/internal/wallet/mq"
	"playground/internal/wallet/usecase"
)

func newTestServer(t *testing.T, r wallet.Repository) *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}
	tm, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
	if err != nil {
		log.Fatal(err)
	}
	p := mq.NewAsynqProducer(asynq.RedisClientOpt{
		Addr: cfg.RedisAddress,
	})
	// usecases
	u := usecase.NewUsecase(r, p, nil, tm, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)
	// handlers
	svr := NewHandler(u, authMiddleware(tm))
	require.NoError(t, err)

	return &Server{
		server: &http.Server{
			Addr:    cfg.HTTPServerAddress,
			Handler: svr,
		},
		tm: tm,
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
