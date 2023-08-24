package usecase

import (
	"context"
	"errors"

	"playground/app"
)

type WalletUsecase struct {
	r app.WalletRepository
}

var _ app.WalletUsecase = (*WalletUsecase)(nil)

func NewWalletUsecase(r app.WalletRepository) app.WalletUsecase {
	return &WalletUsecase{
		r: r,
	}
}

func (u *WalletUsecase) Create(ctx context.Context, args app.CreateWalletInputParams) (*app.Wallet, error) {
	return u.r.Create(ctx, app.NewWallet(args.UserID, args.Balance, args.Currency))
}

func (u *WalletUsecase) Get(ctx context.Context, args app.GetOrDeleteWalletInputParams) (*app.Wallet, error) {
	w, err := u.r.Get(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if w.UserID != args.UserID {
		return nil, errors.New("wallet doesn't belong to the authenticated user")
	}

	return w, nil
}

func (u *WalletUsecase) List(ctx context.Context, arg app.ListWalletsInputParams) ([]app.Wallet, error) {
	return u.r.List(ctx, arg)
}

func (u *WalletUsecase) Delete(ctx context.Context, args app.GetOrDeleteWalletInputParams) error {
	w, err := u.r.Get(ctx, args.ID)
	if err != nil {
		return err
	}
	if w.UserID != args.UserID {
		return errors.New("wallet doesn't belong to the authenticated user")
	}

	return u.r.Delete(ctx, args.ID)
}
