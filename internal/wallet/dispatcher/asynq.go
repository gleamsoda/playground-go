package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"

	"playground/internal/wallet"
)

type Dispatcher struct {
	c *asynq.Client
}

func NewDispatcher(i *do.Injector) (wallet.Dispatcher, error) {
	redisOpt := do.MustInvoke[asynq.RedisClientOpt](i)
	return &Dispatcher{
		c: asynq.NewClient(redisOpt),
	}, nil
}

func (p *Dispatcher) SendVerifyEmail(ctx context.Context, payload *wallet.SendVerifyEmailPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
	}
	task := asynq.NewTask(wallet.SendVerifyEmailQueue, jsonPayload, opts...)
	info, err := p.c.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}
