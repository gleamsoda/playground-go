package usecase

import (
	"context"

	"playground/app"
)

type AccountUsecase struct {
	ar app.AccountRepository
}

func NewAccountUsecase(ar app.AccountRepository) app.AccountUsecase {
	return &AccountUsecase{
		ar: ar,
	}
}

func (u *AccountUsecase) CreateAccount(ctx context.Context, args *app.CreateAccountParams) (*app.Account, error) {
	return u.ar.CreateAccount(ctx, app.NewAccount(args.Owner, args.Balance, args.Currency))
}

func (u *AccountUsecase) GetAccount(ctx context.Context, id int64) (*app.Account, error) {
	return u.ar.GetAccount(ctx, id)
}

func (u *AccountUsecase) ListAccounts(ctx context.Context, args *app.ListAccountsParams) ([]app.Account, error) {
	return u.ar.ListAccounts(ctx, args)
}
