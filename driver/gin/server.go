package gin

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"

	"playground/app/mq"
	"playground/app/repository"
	"playground/app/usecase"
	"playground/config"
	"playground/pkg/token"
)

type Server struct {
	server *http.Server
	tm     token.Manager
}

func NewServer(cfg config.Config) (*Server, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}
	tm, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	// repositories
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/playground?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort))
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	p := mq.NewAsynqProducer(asynq.RedisClientOpt{
		Addr: cfg.RedisAddress,
	})
	r := repository.NewRepository(conn)
	// usecases
	u := usecase.NewUsecase(r, p, nil, tm, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)
	// handlers
	svr := NewHandler(u, authMiddleware(tm))

	return &Server{
		server: &http.Server{
			Addr:    cfg.HTTPServerAddress,
			Handler: svr,
		},
		tm: tm,
	}, nil
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
