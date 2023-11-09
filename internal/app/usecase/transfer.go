package usecase

import (
	"context"
	"errors"
	"fmt"

	"playground/internal/app"
)

type CreateTransfer struct {
	r app.RepositoryManager
}

func NewCreateTransfer(r app.RepositoryManager) *CreateTransfer {
	return &CreateTransfer{
		r: r,
	}
}

func (u *CreateTransfer) Execute(ctx context.Context, args *app.CreateTransferParams) (*app.Transfer, error) {
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

	return u.r.Transfer().Create(ctx, app.NewTransfer(args.FromAccountID, args.ToAccountID, args.Amount))
}

func (u *CreateTransfer) validAccount(ctx context.Context, accountID int64, currency string) (*app.Account, error) {
	a, err := u.r.Account().Get(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if a.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", a.ID, a.Currency, currency)
		return nil, err
	}

	return a, nil
}
