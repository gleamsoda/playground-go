package asynq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"

	"playground/internal/app"
	"playground/internal/app/usecase"
	"playground/internal/pkg/mail"
)

type handler struct {
	sendVerifyEmail app.SendVerifyEmailUsecase
}

func NewHandler(i *do.Injector) (*handler, error) {
	r := do.MustInvoke[app.RepositoryManager](i)
	m := do.MustInvoke[mail.Sender](i)
	return &handler{
		sendVerifyEmail: usecase.NewSendVerifyEmail(r, m),
	}, nil
}

func (h *handler) SendVerifyEmail(ctx context.Context, t *asynq.Task) error {
	var payload app.SendVerifyEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	usr, err := h.sendVerifyEmail.Execute(ctx, &app.SendVerifyEmailPayload{
		Username: payload.Username,
	})
	if err != nil {
		return err
	}

	log.Info().Str("type", t.Type()).Bytes("payload", t.Payload()).
		Str("email", usr.Email).Msg("processed task")
	return nil
}
