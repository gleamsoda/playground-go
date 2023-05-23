package usecase

import (
	"context"

	"github.com/gleamsoda/go-playground/domain"
	"github.com/gleamsoda/go-playground/internal"
)

type UserUsecase struct {
	r domain.UserRepository
}

var _ domain.UserUsecase = (*UserUsecase)(nil)

func NewUserUsecase(r domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		r: r,
	}
}

func (u *UserUsecase) Create(ctx context.Context, arg domain.CreateUserParams) (*domain.User, error) {
	hashedPassword, err := internal.HashPassword(arg.Password)
	if err != nil {
		return nil, err
	}
	usr := domain.NewUser(
		arg.Username,
		arg.FullName,
		arg.Email,
		hashedPassword,
	)
	return u.r.Create(ctx, usr)
}

func (u *UserUsecase) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return u.r.GetByUsername(ctx, username)
}
