package asynq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"playground/internal/wallet"
)

type handler struct {
	w wallet.Usecase
}

func NewHandler(w wallet.Usecase) *handler {
	return &handler{w: w}
}

func (h *handler) SendVerifyEmail(ctx context.Context, t *asynq.Task) error {
	var payload wallet.SendVerifyEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	usr, err := h.w.SendVerifyEmail(ctx, &wallet.SendVerifyEmailPayload{
		Username: payload.Username,
	})
	if err != nil {
		return err
	}

	log.Info().Str("type", t.Type()).Bytes("payload", t.Payload()).
		Str("email", usr.Email).Msg("processed task")
	return nil
}
