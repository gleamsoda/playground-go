package usecase

import (
	"context"
	"errors"

	"github.com/morikuni/failure"

	"playground/app"
	"playground/pkg/apperrors"
)

func (u *Usecase) CreateAccount(ctx context.Context, args *app.CreateAccountParams) (*app.Account, error) {
	return u.r.CreateAccount(ctx, app.NewAccount(args.Owner, args.Balance, args.Currency))
}

func (u *Usecase) GetAccount(ctx context.Context, args *app.GetAccountsParams) (*app.Account, error) {
	a, err := u.r.GetAccount(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if a.Owner != args.Owner {
		err := errors.New("account doesn't belong to the authenticated user")
		return nil, failure.Translate(err, apperrors.Unauthorized)
	}
	return a, nil
}

func (u *Usecase) ListAccounts(ctx context.Context, args *app.ListAccountsParams) ([]app.Account, error) {
	return u.r.ListAccounts(ctx, args)
}
