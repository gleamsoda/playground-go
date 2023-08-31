package usecase

import (
	"context"
	"errors"
	"fmt"

	"playground/app"
)

type TransferUsecase struct {
	ar app.AccountRepository
	tr app.TransferRepository
}

func NewTransferUsecase(tr app.TransferRepository, ar app.AccountRepository) app.TransferUsecase {
	return &TransferUsecase{
		tr: tr,
		ar: ar,
	}
}

func (u *TransferUsecase) CreateTransfer(ctx context.Context, args *app.CreateTransferParams) (*app.Transfer, error) {
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

	return u.tr.CreateTransfer(ctx, app.NewTransfer(args.FromAccountID, args.ToAccountID, args.Amount))
}

func (u *TransferUsecase) validAccount(ctx context.Context, accountID int64, currency string) (*app.Account, error) {
	a, err := u.ar.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if a.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", a.ID, a.Currency, currency)
		return nil, err
	}

	return a, nil
}
