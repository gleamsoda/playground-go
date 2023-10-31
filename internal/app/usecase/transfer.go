package usecase

import (
	"context"
	"errors"
	"fmt"

	"playground/internal/app"
)

type CreateTransferUsecase struct {
	r app.Repository
}

func NewCreateTransferUsecase(r app.Repository) *CreateTransferUsecase {
	return &CreateTransferUsecase{
		r: r,
	}
}

func (u *CreateTransferUsecase) Execute(ctx context.Context, args *app.CreateTransferParams) (*app.Transfer, error) {
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

	return u.r.CreateTransfer(ctx, app.NewTransfer(args.FromAccountID, args.ToAccountID, args.Amount))
}

func (u *CreateTransferUsecase) validAccount(ctx context.Context, accountID int64, currency string) (*app.Account, error) {
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
