package usecase

import (
	"context"

	"github.com/gleamsoda/go-playground/domain"
)

type TransferUsecase struct {
	transferRepo domain.TransferRepository
	entryRepo    domain.EntryRepository
	walletRepo   domain.WalletRepository
}

var _ domain.TransferUsecase = (*TransferUsecase)(nil)

func NewTransferUsecase(transferRepo domain.TransferRepository, entryRepo domain.EntryRepository, walletRepo domain.WalletRepository) *TransferUsecase {
	return &TransferUsecase{
		transferRepo: transferRepo,
		entryRepo:    entryRepo,
		walletRepo:   walletRepo,
	}
}

func (u *TransferUsecase) Create(ctx context.Context, arg domain.CreateTransferParams) (*domain.Transfer, error) {
	return u.transferRepo.Create(ctx, arg)
}
