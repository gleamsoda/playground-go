package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"playground/domain"
	"playground/repo/sqlc/gen"
)

type SessionRepository struct {
	q *gen.Queries
}

var _ domain.SessionRepository = (*SessionRepository)(nil)

func NewSessionRepository(db *sql.DB) domain.SessionRepository {
	return &SessionRepository{
		q: gen.New(db),
	}
}

func (r *SessionRepository) Create(ctx context.Context, arg *domain.Session) error {
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

func (r *SessionRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	session, err := r.q.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Session{
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
