package usecase

import (
	"context"
	"errors"

	"github.com/morikuni/failure"

	"playground/internal/app"
	"playground/internal/pkg/apperr"
)

type (
	CreateAccountUsecase struct {
		r app.Repository
	}
	GetAccountUsecase struct {
		r app.Repository
	}
	ListAccountsUsecase struct {
		r app.Repository
	}
)

func NewCreateAccountUsecase(r app.Repository) *CreateAccountUsecase {
	return &CreateAccountUsecase{
		r: r,
	}
}

func (u *CreateAccountUsecase) Execute(ctx context.Context, args *app.CreateAccountParams) (*app.Account, error) {
	return u.r.CreateAccount(ctx, app.NewAccount(args.Owner, args.Balance, args.Currency))
}

func NewGetAccountUsecase(r app.Repository) *GetAccountUsecase {
	return &GetAccountUsecase{
		r: r,
	}
}

func (u *GetAccountUsecase) Execute(ctx context.Context, args *app.GetAccountsParams) (*app.Account, error) {
	a, err := u.r.GetAccount(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if a.Owner != args.Owner {
		err := errors.New("account doesn't belong to the authenticated user")
		return nil, failure.Translate(err, apperr.Unauthenticated)
	}
	return a, nil
}

func NewListAccountsUsecase(r app.Repository) *ListAccountsUsecase {
	return &ListAccountsUsecase{
		r: r,
	}
}

func (u *ListAccountsUsecase) Execute(ctx context.Context, args *app.ListAccountsParams) ([]app.Account, error) {
	return u.r.ListAccounts(ctx, args)
}
