package usecase

import (
	"context"
	"errors"

	"playground/app"
)

type EntryUsecase struct {
	entryRepo  app.EntryRepository
	walletRepo app.WalletRepository
}

var _ app.EntryUsecase = (*EntryUsecase)(nil)

func NewEntryUsecase(entryRepo app.EntryRepository, walletRepo app.WalletRepository) app.EntryUsecase {
	return &EntryUsecase{
		entryRepo:  entryRepo,
		walletRepo: walletRepo,
	}
}

func (u *EntryUsecase) List(ctx context.Context, arg app.ListEntriesInputParams) ([]app.Entry, error) {
	w, err := u.walletRepo.Get(ctx, arg.WalletID)
	if err != nil {
		return nil, err
	}

	if w.UserID != arg.RequestUserID {
		return nil, errors.New("wallet doesn't belong to the authenticated user")
	}
	return u.entryRepo.List(ctx, arg)
}
