package usecase

import (
	"context"
	"errors"

	"github.com/morikuni/failure"

	"playground/internal/pkg/apperr"
	"playground/internal/wallet"
)

func (u *Usecase) CreateAccount(ctx context.Context, args *wallet.CreateAccountParams) (*wallet.Account, error) {
	return u.r.CreateAccount(ctx, wallet.NewAccount(args.Owner, args.Balance, args.Currency))
}

func (u *Usecase) GetAccount(ctx context.Context, args *wallet.GetAccountsParams) (*wallet.Account, error) {
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

func (u *Usecase) ListAccounts(ctx context.Context, args *wallet.ListAccountsParams) ([]wallet.Account, error) {
	return u.r.ListAccounts(ctx, args)
}
