package asynq

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"

	"playground/internal/app"
	"playground/internal/app/repository"
	"playground/internal/config"
	"playground/internal/pkg/mail"
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
	mailer := mail.NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)

	injector := do.New()
	do.Provide(injector, NewHandler)
	do.Provide(injector, repository.NewManager)
	do.ProvideValue(injector, conn)
	do.ProvideValue(injector, mailer)
	h := do.MustInvoke[*handler](injector)

	return &Consumer{
		server:  s,
		handler: h,
	}, nil
}

func (c *Consumer) Run() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(app.SendVerifyEmailQueue, c.handler.SendVerifyEmail)
	return c.server.Run(mux)
}
