package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/morikuni/failure"

	"playground/internal/app"
	"playground/internal/app/repository/sqlc/gen"
	"playground/internal/pkg/apperr"
)

type User struct {
	exec Executor
	q    gen.Querier
}

func NewUser(e Executor) *User {
	return &User{
		exec: e,
		q:    gen.New(e),
	}
}

var _ app.UserRepository = (*User)(nil)

func (r *User) Create(ctx context.Context, args *app.User) (*app.User, error) {
	if _, err := r.q.CreateUser(ctx, &gen.CreateUserParams{
		Username:       args.Username,
		FullName:       args.FullName,
		Email:          args.Email,
		HashedPassword: args.HashedPassword,
	}); err != nil {
		return nil, err
	}

	return r.Get(ctx, args.Username)
}

func (r *User) Get(ctx context.Context, username string) (*app.User, error) {
	u, err := r.q.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.Translate(err, apperr.NotFound)
		}
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

func (r *User) Update(ctx context.Context, args *app.User) (*app.User, error) {
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

	return r.Get(ctx, args.Username)
}

func (r *User) UpdateEmailVerified(ctx context.Context, args *app.VerifyEmailParams) (*app.User, *app.VerifyEmail, error) {
	var u *gen.User
	var ve *gen.VerifyEmail
	if err := runTx(ctx, r.exec, func(ctx context.Context, tx *sql.Tx) error {
		txr := NewUser(tx)
		var err error
		if err = txr.q.UpdateVerifyEmail(ctx, &gen.UpdateVerifyEmailParams{
			ID:         args.EmailID,
			SecretCode: args.SecretCode,
		}); err != nil {
			return err
		}
		ve, err = txr.q.GetVerifyEmail(ctx, args.EmailID)
		if err != nil {
			if err == sql.ErrNoRows {
				return failure.Translate(err, apperr.NotFound)
			}
			return err
		}
		u, err = txr.q.GetUser(ctx, ve.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return failure.Translate(err, apperr.NotFound)
			}
			return err
		}
		u.IsEmailVerified = true
		if err := txr.q.UpdateUser(ctx, &gen.UpdateUserParams{
			HashedPassword:    u.HashedPassword,
			PasswordChangedAt: u.PasswordChangedAt,
			FullName:          u.FullName,
			Email:             u.Email,
			IsEmailVerified:   u.IsEmailVerified,
			Username:          u.Username,
		}); err != nil {
			return err
		}
		u, err = txr.q.GetUser(ctx, u.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return failure.Translate(err, apperr.NotFound)
			}
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}

	return &app.User{
			Username:          u.Username,
			HashedPassword:    u.HashedPassword,
			FullName:          u.FullName,
			Email:             u.Email,
			PasswordChangedAt: u.PasswordChangedAt,
			CreatedAt:         u.CreatedAt,
			IsEmailVerified:   u.IsEmailVerified,
		}, &app.VerifyEmail{
			ID:         ve.ID,
			Username:   ve.Username,
			Email:      ve.Email,
			SecretCode: ve.SecretCode,
			IsUsed:     ve.IsUsed,
			ExpiredAt:  ve.ExpiredAt,
			CreatedAt:  ve.CreatedAt,
		}, nil
}

func (r *User) CreateSession(ctx context.Context, args *app.Session) error {
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

func (r *User) GetSession(ctx context.Context, id uuid.UUID) (*app.Session, error) {
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

func (r *User) GetVerifyEmail(ctx context.Context, id int64) (*app.VerifyEmail, error) {
	ve, err := r.q.GetVerifyEmail(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, failure.Translate(err, apperr.NotFound)
		}
		return nil, err
	}

	return &app.VerifyEmail{
		ID:         ve.ID,
		Username:   ve.Username,
		Email:      ve.Email,
		SecretCode: ve.SecretCode,
		IsUsed:     ve.IsUsed,
		ExpiredAt:  ve.ExpiredAt,
		CreatedAt:  ve.CreatedAt,
	}, nil
}

func (r *User) CreateVerifyEmail(ctx context.Context, args *app.VerifyEmail) (*app.VerifyEmail, error) {
	id, err := r.q.CreateVerifyEmail(ctx, &gen.CreateVerifyEmailParams{
		Username:   args.Username,
		Email:      args.Email,
		SecretCode: args.SecretCode,
	})
	if err != nil {
		return nil, err
	}

	return r.GetVerifyEmail(ctx, id)
}
