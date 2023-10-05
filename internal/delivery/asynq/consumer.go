package asynq

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"playground/internal/config"
	"playground/internal/pkg/mail"
	"playground/internal/pkg/token"
	"playground/internal/wallet"
	"playground/internal/wallet/repository"
	"playground/internal/wallet/usecase"
)

type Consumer struct {
	server  *asynq.Server
	handler *handler
}

func Run() error {
	if c, err := NewConsumer(config.Get()); err != nil {
		return err
	} else {
		return c.Run()
	}
}

func NewConsumer(cfg config.Config) (*Consumer, error) {
	redis.SetLogger(lgr)
	s := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: cfg.RedisAddress,
		},
		asynq.Config{
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, t *asynq.Task, err error) {
				log.Error().Err(err).Str("type", t.Type()).
					Bytes("payload", t.Payload()).Msg("process task failed")
			}),
			Logger: lgr,
		},
	)

	conn, err := sql.Open("mysql", cfg.DBName())
	if err != nil {
		return nil, err
	} else if err := conn.Ping(); err != nil {
		return nil, err
	}
	tm, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	r := repository.NewRepository(conn)
	mailer := mail.NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)
	u := usecase.NewUsecase(r, nil, mailer, tm, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)

	return &Consumer{
		server:  s,
		handler: NewHandler(u),
	}, nil
}

func (c *Consumer) Run() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(wallet.SendVerifyEmailQueue, c.handler.SendVerifyEmail)
	return c.server.Run(mux)
}
