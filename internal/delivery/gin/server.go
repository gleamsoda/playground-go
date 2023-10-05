package gin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"playground/internal/config"
	"playground/internal/delivery/gin/handler"
	"playground/internal/delivery/gin/middleware"
	"playground/internal/pkg/token"
	"playground/internal/wallet"
	"playground/internal/wallet/dispatcher"
	"playground/internal/wallet/repository"
	"playground/internal/wallet/usecase"
)

type Server struct {
	server *http.Server
	tm     token.Manager
}

func Run(ctx context.Context) error {
	server, err := NewServer(config.Get())
	if err != nil {
		return err
	}

	ctxn, canceln := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer canceln()
	ctxs, cancels := context.WithCancelCause(ctx)
	defer cancels(nil)
	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			cancels(err)
		}
		cancels(nil)
	}()

	select {
	case <-ctxn.Done():
		log.Info().Msg("shutting down gracefully...")
		ctxsd, cancelsd := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancelsd()
		if err := server.Shutdown(ctxsd); err != nil {
			return fmt.Errorf("failed to shutdown gracefully: %v", err)
		}
		log.Info().Msg("shutdown complete")
	case <-ctxs.Done():
		if cause := context.Cause(ctxs); !errors.Is(cause, context.Canceled) {
			return fmt.Errorf("error: %v", cause)
		}
	}
	return nil
}

func NewServer(cfg config.Config) (*Server, error) {
	// validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}
	tm, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	// repositories
	conn, err := sql.Open("mysql", cfg.DBName())
	if err != nil {
		return nil, err
	} else if err := conn.Ping(); err != nil {
		return nil, err
	}
	r := repository.NewRepository(conn)
	p := dispatcher.NewDispatcher(asynq.RedisClientOpt{
		Addr: cfg.RedisAddress,
	})
	// usecases
	u := usecase.NewUsecase(r, p, nil, tm, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)

	// handlers
	router := NewRouter(handler.NewHandler(u), middleware.Auth(tm))
	return &Server{
		server: &http.Server{
			Addr:    cfg.HTTPServerAddress,
			Handler: router,
		},
		tm: tm,
	}, nil
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if c, ok := fieldLevel.Field().Interface().(string); ok {
		return wallet.IsSupportedCurrency(c)
	}
	return false
}
