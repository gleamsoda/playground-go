package usecase

import (
	"context"

	"playground/app"
)

func (u *Usecase) CreateAccount(ctx context.Context, args *app.CreateAccountParams) (*app.Account, error) {
	return u.r.CreateAccount(ctx, app.NewAccount(args.Owner, args.Balance, args.Currency))
}

func (u *Usecase) GetAccount(ctx context.Context, id int64) (*app.Account, error) {
	return u.r.GetAccount(ctx, id)
}

func (u *Usecase) ListAccounts(ctx context.Context, args *app.ListAccountsParams) ([]app.Account, error) {
	return u.r.ListAccounts(ctx, args)
}
