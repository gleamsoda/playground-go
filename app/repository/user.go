package repository

import (
	"context"

	"github.com/google/uuid"

	"playground/app"
	"playground/app/repository/sqlc/gen"
)

func (r *Repository) CreateUser(ctx context.Context, args *app.User) (*app.User, error) {
	if _, err := r.q.CreateUser(ctx, &gen.CreateUserParams{
		Username:       args.Username,
		FullName:       args.FullName,
		Email:          args.Email,
		HashedPassword: args.HashedPassword,
	}); err != nil {
		return nil, err
	}

	return r.GetUser(ctx, args.Username)
}

func (r *Repository) GetUser(ctx context.Context, username string) (*app.User, error) {
	u, err := r.q.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	return &app.User{
		Username:          u.Username,
		HashedPassword:    u.HashedPassword,
		FullName:          u.FullName,
		Email:             u.Email,
		PasswordChangedAt: u.PasswordChangedAt,
		CreatedAt:         u.CreatedAt,
	}, nil
}

func (r *Repository) CreateSession(ctx context.Context, args *app.Session) error {
	return r.q.CreateSession(ctx, &gen.CreateSessionParams{
		ID:           args.ID,
		Username:     args.Username,
		RefreshToken: args.RefreshToken,
		UserAgent:    args.UserAgent,
		ClientIp:     args.ClientIP,
		IsBlocked:    args.IsBlocked,
		ExpiresAt:    args.ExpiresAt,
	})
}

func (r *Repository) GetSession(ctx context.Context, id uuid.UUID) (*app.Session, error) {
	s, err := r.q.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	return &app.Session{
		ID:           s.ID,
		Username:     s.Username,
		RefreshToken: s.RefreshToken,
		UserAgent:    s.UserAgent,
		ClientIP:     s.ClientIp,
		IsBlocked:    s.IsBlocked,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
	}, nil
}
