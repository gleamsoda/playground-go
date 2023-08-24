package repository

import (
	"context"
	"database/sql"

	"playground/app"
	"playground/app/repository/sqlc/gen"
)

type UserRepository struct {
	q *gen.Queries
}

var _ app.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) app.UserRepository {
	return &UserRepository{
		q: gen.New(db),
	}
}

func (r *UserRepository) Create(ctx context.Context, arg *app.User) (*app.User, error) {
	_, err := r.q.CreateUser(ctx, gen.CreateUserParams{
		Username:       arg.Username,
		FullName:       arg.FullName,
		Email:          arg.Email,
		HashedPassword: arg.HashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return r.GetByUsername(ctx, arg.Username)
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*app.User, error) {
	user, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &app.User{
		ID:             user.ID,
		Username:       user.Username,
		FullName:       user.FullName,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		CreatedAt:      user.CreatedAt,
	}, nil
}
