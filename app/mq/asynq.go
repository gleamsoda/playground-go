package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type AsynqProducer struct {
	c *asynq.Client
}

func NewAsynqProducer(redisOpt asynq.RedisClientOpt) Producer {
	return &AsynqProducer{
		c: asynq.NewClient(redisOpt),
	}
}

func (p *AsynqProducer) SendVerifyEmail(ctx context.Context, payload *SendVerifyEmailPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
	}
	task := asynq.NewTask(SendVerifyEmailQueue, jsonPayload, opts...)
	info, err := p.c.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}
