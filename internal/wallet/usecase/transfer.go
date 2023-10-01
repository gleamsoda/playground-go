package usecase

import (
	"context"
	"errors"
	"fmt"

	"playground/internal/wallet"
)

func (u *Usecase) CreateTransfer(ctx context.Context, args *wallet.CreateTransferParams) (*wallet.Transfer, error) {
	fromAccount, err := u.validAccount(ctx, args.FromAccountID, args.Currency)
	if err != nil {
		return nil, err
	}
	if fromAccount.Owner != args.RequestUsername {
		err := errors.New("from account doesn't belong to the authenticated user")
		return nil, err
	}
	if _, err = u.validAccount(ctx, args.ToAccountID, args.Currency); err != nil {
		return nil, err
	}

	return u.r.CreateTransfer(ctx, wallet.NewTransfer(args.FromAccountID, args.ToAccountID, args.Amount))
}

func (u *Usecase) validAccount(ctx context.Context, accountID int64, currency string) (*wallet.Account, error) {
	a, err := u.r.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if a.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", a.ID, a.Currency, currency)
		return nil, err
	}

	return a, nil
}
