package asynq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"

	"playground/internal/app"
)

type handler struct {
	w app.Usecase
}

func NewHandler(i *do.Injector) (*handler, error) {
	w := do.MustInvoke[app.Usecase](i)
	return &handler{w: w}, nil
}

func (h *handler) SendVerifyEmail(ctx context.Context, t *asynq.Task) error {
	var payload app.SendVerifyEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	usr, err := h.w.SendVerifyEmail(ctx, &app.SendVerifyEmailPayload{
		Username: payload.Username,
	})
	if err != nil {
		return err
	}

	log.Info().Str("type", t.Type()).Bytes("payload", t.Payload()).
		Str("email", usr.Email).Msg("processed task")
	return nil
}
