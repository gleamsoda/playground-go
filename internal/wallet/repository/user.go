package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/morikuni/failure"

	"playground/internal/pkg/apperr"
	"playground/internal/wallet"
	"playground/internal/wallet/repository/sqlc/gen"
)

func (r *Repository) CreateUser(ctx context.Context, args *wallet.User) (*wallet.User, error) {
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

func (r *Repository) GetUser(ctx context.Context, username string) (*wallet.User, error) {
	u, err := r.q.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.Translate(err, apperr.NotFound)
		}
		return nil, err
	}

	return &wallet.User{
		Username:          u.Username,
		HashedPassword:    u.HashedPassword,
		FullName:          u.FullName,
		Email:             u.Email,
		PasswordChangedAt: u.PasswordChangedAt,
		CreatedAt:         u.CreatedAt,
	}, nil
}

func (r *Repository) UpdateUser(ctx context.Context, args *wallet.User) (*wallet.User, error) {
	err := r.q.UpdateUser(ctx, &gen.UpdateUserParams{
		Username:          args.Username,
		FullName:          args.FullName,
		Email:             args.Email,
		HashedPassword:    args.HashedPassword,
		PasswordChangedAt: args.PasswordChangedAt,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.Translate(err, apperr.NotFound)
		}
		return nil, err
	}

	return r.GetUser(ctx, args.Username)
}

func (r *Repository) CreateSession(ctx context.Context, args *wallet.Session) error {
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

func (r *Repository) GetSession(ctx context.Context, id uuid.UUID) (*wallet.Session, error) {
	s, err := r.q.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	return &wallet.Session{
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
