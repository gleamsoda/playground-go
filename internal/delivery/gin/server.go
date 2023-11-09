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

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"

	"playground/internal/app"
	"playground/internal/app/dispatcher"
	"playground/internal/app/repository"
	"playground/internal/config"
	"playground/internal/delivery/gin/handler"
	"playground/internal/pkg/token"
)

type Server struct {
	server *http.Server
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
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}
	tm, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	conn, err := sql.Open("mysql", cfg.DBName())
	if err != nil {
		return nil, err
	} else if err := conn.Ping(); err != nil {
		return nil, err
	}

	injector := do.New()
	do.Provide(injector, NewRouter)
	do.Provide(injector, handler.NewHandler)
	do.Provide(injector, repository.NewManager)
	do.ProvideValue(injector, conn)
	do.Provide(injector, dispatcher.NewDispatcher)
	do.ProvideValue(injector, asynq.RedisClientOpt{Addr: cfg.RedisAddress})
	do.ProvideValue(injector, tm)
	do.ProvideNamedValue(injector, "AccessTokenDuration", cfg.AccessTokenDuration)
	do.ProvideNamedValue(injector, "RefreshTokenDuration", cfg.RefreshTokenDuration)

	// handlers
	router := do.MustInvoke[*gin.Engine](injector)
	return &Server{
		server: &http.Server{
			Addr:    cfg.HTTPServerAddress,
			Handler: router,
		},
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
		return app.IsSupportedCurrency(c)
	}
	return false
}
