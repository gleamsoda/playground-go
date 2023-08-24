package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"playground/app"
	"playground/app/repository/gen"
)

type SessionRepository struct {
	q *gen.Queries
}

var _ app.SessionRepository = (*SessionRepository)(nil)

func NewSessionRepository(db *sql.DB) app.SessionRepository {
	return &SessionRepository{
		q: gen.New(db),
	}
}

func (r *SessionRepository) Create(ctx context.Context, arg *app.Session) error {
	return r.q.CreateSession(ctx, gen.CreateSessionParams{
		ID:           arg.ID,
		UserID:       arg.UserID,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIP,
		IsBlocked:    arg.IsBlocked,
		ExpiresAt:    arg.ExpiresAt,
	})
}

func (r *SessionRepository) Get(ctx context.Context, id uuid.UUID) (*app.Session, error) {
	session, err := r.q.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	return &app.Session{
		ID:           session.ID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIP:     session.ClientIp,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}, nil
}
