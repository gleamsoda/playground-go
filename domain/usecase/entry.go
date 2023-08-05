package usecase

import (
	"context"
	"errors"

	"playground/domain"
)

type EntryUsecase struct {
	entryRepo  domain.EntryRepository
	walletRepo domain.WalletRepository
}

var _ domain.EntryUsecase = (*EntryUsecase)(nil)

func NewEntryUsecase(entryRepo domain.EntryRepository, walletRepo domain.WalletRepository) domain.EntryUsecase {
	return &EntryUsecase{
		entryRepo:  entryRepo,
		walletRepo: walletRepo,
	}
}

func (u *EntryUsecase) List(ctx context.Context, arg domain.ListEntriesInputParams) ([]domain.Entry, error) {
	w, err := u.walletRepo.Get(ctx, arg.WalletID)
	if err != nil {
		return nil, err
	}

	if w.UserID != arg.RequestUserID {
		return nil, errors.New("wallet doesn't belong to the authenticated user")
	}
	return u.entryRepo.List(ctx, arg)
}
