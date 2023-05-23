package usecase

import (
	"context"

	"github.com/gleamsoda/go-playground/domain"
)

type WalletUsecase struct {
	r domain.WalletRepository
}

var _ domain.WalletUsecase = (*WalletUsecase)(nil)

func NewWalletUsecase(r domain.WalletRepository) *WalletUsecase {
	return &WalletUsecase{
		r: r,
	}
}

func (u *WalletUsecase) Create(ctx context.Context, args domain.CreateWalletParams) (*domain.Wallet, error) {
	return u.r.Create(ctx, args)
}

func (u *WalletUsecase) Get(ctx context.Context, id int64) (*domain.Wallet, error) {
	return u.r.Get(ctx, id)
}

func (u *WalletUsecase) List(ctx context.Context, arg domain.ListWalletsParams) ([]domain.Wallet, error) {
	return u.r.List(ctx, arg)
}

func (u *WalletUsecase) Delete(ctx context.Context, id int64) error {
	return u.r.Delete(ctx, id)
}
