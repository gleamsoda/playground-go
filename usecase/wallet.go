package usecase

import (
	"context"
	"errors"

	"playground/domain"
)

type WalletUsecase struct {
	r domain.WalletRepository
}

var _ domain.WalletUsecase = (*WalletUsecase)(nil)

func NewWalletUsecase(r domain.WalletRepository) domain.WalletUsecase {
	return &WalletUsecase{
		r: r,
	}
}

func (u *WalletUsecase) Create(ctx context.Context, args domain.CreateWalletInputParams) (*domain.Wallet, error) {
	return u.r.Create(ctx, domain.NewWallet(args.UserID, args.Balance, args.Currency))
}

func (u *WalletUsecase) Get(ctx context.Context, args domain.GetOrDeleteWalletInputParams) (*domain.Wallet, error) {
	w, err := u.r.Get(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if w.UserID != args.UserID {
		return nil, errors.New("wallet doesn't belong to the authenticated user")
	}

	return w, nil
}

func (u *WalletUsecase) List(ctx context.Context, arg domain.ListWalletsInputParams) ([]domain.Wallet, error) {
	return u.r.List(ctx, arg)
}

func (u *WalletUsecase) Delete(ctx context.Context, args domain.GetOrDeleteWalletInputParams) error {
	w, err := u.r.Get(ctx, args.ID)
	if err != nil {
		return err
	}
	if w.UserID != args.UserID {
		return errors.New("wallet doesn't belong to the authenticated user")
	}

	return u.r.Delete(ctx, args.ID)
}
