package asynq

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"playground/app"
	"playground/app/mq"
	"playground/app/repository"
	"playground/app/usecase"
	"playground/config"
	"playground/internal/pkg/mail"
	"playground/internal/pkg/token"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type Consumer struct {
	s *asynq.Server
	u app.Usecase
}

func NewConsumer(cfg config.Config) (*Consumer, error) {
	logger := NewLogger()
	redis.SetLogger(logger)

	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: cfg.RedisAddress,
		},
		asynq.Config{
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: logger,
		},
	)

	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/playground?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort))
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
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
		s: server,
		u: u,
	}, nil
}

func (c *Consumer) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(mq.SendVerifyEmailQueue, c.SendVerifyEmail)
	return c.s.Run(mux)
}

func (c *Consumer) Stop() {
	c.s.Shutdown()
}

func (c *Consumer) SendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload mq.SendVerifyEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	usr, err := c.u.SendVerifyEmail(ctx, &mq.SendVerifyEmailPayload{
		Username: payload.Username,
	})
	if err != nil {
		return err
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", usr.Email).Msg("processed task")
	return nil
}
