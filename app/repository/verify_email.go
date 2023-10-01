package repository

import (
	"context"
	"database/sql"

	"github.com/morikuni/failure"

	"playground/app"
	"playground/app/repository/sqlc/gen"
	"playground/pkg/apperrors"
)

func (r *Repository) GetVerifyEmail(ctx context.Context, id int64) (*app.VerifyEmail, error) {
	ve, err := r.q.GetVerifyEmail(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, failure.Translate(err, apperrors.NotFound)
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

func (r *Repository) CreateVerifyEmail(ctx context.Context, args *app.VerifyEmail) (*app.VerifyEmail, error) {
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

func (r *Repository) UpdateVerifyEmail(ctx context.Context, args *app.VerifyEmail) (*app.VerifyEmail, error) {
	err := r.q.UpdateVerifyEmail(ctx, &gen.UpdateVerifyEmailParams{
		ID:         args.ID,
		SecretCode: args.SecretCode,
	})
	if err != nil {
		return nil, err
	}

	return r.GetVerifyEmail(ctx, args.ID)
}

func (r *Repository) UpdateUserEmailVerified(ctx context.Context, args *app.VerifyEmailParams) (*app.User, *app.VerifyEmail, error) {
	var u *gen.User
	var ve *gen.VerifyEmail
	if err := r.tx(ctx, func(ctx context.Context, q *gen.Queries) error {
		var err error
		if err = q.UpdateVerifyEmail(ctx, &gen.UpdateVerifyEmailParams{
			ID:         args.EmailID,
			SecretCode: args.SecretCode,
		}); err != nil {
			return err
		}
		ve, err = q.GetVerifyEmail(ctx, args.EmailID)
		if err != nil {
			if err == sql.ErrNoRows {
				return failure.Translate(err, apperrors.NotFound)
			}
			return err
		}
		u, err = q.GetUser(ctx, ve.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return failure.Translate(err, apperrors.NotFound)
			}
			return err
		}
		u.IsEmailVerified = true
		if err := q.UpdateUser(ctx, &gen.UpdateUserParams{
			HashedPassword:    u.HashedPassword,
			PasswordChangedAt: u.PasswordChangedAt,
			FullName:          u.FullName,
			Email:             u.Email,
			IsEmailVerified:   u.IsEmailVerified,
			Username:          u.Username,
		}); err != nil {
			return err
		}
		u, err = q.GetUser(ctx, u.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return failure.Translate(err, apperrors.NotFound)
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
