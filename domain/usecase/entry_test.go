package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"playground/domain"
	mock_domain "playground/test/mock/domain"
)

func TestEntryUsecase_List(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mock_domain.NewMockWalletRepository(ctrl)
	mockEntryRepo := mock_domain.NewMockEntryRepository(ctrl)
	u := &EntryUsecase{
		entryRepo:  mockEntryRepo,
		walletRepo: mockWalletRepo,
	}
	ctx := context.Background()
	args := domain.ListEntriesInputParams{
		WalletID:      123,
		RequestUserID: 123,
	}

	t.Run("Wallet belongs to the user", func(t *testing.T) {
		wallet := &domain.Wallet{UserID: 123}
		es := []domain.Entry{
			{ID: 123, WalletID: 123},
			{ID: 456, WalletID: 123},
		}
		mockWalletRepo.EXPECT().Get(ctx, args.WalletID).Return(wallet, nil)
		mockEntryRepo.EXPECT().List(ctx, args).Return(es, nil)

		entries, err := u.List(ctx, args)
		assert.NoError(t, err)
		assert.Equal(t, []domain.Entry{
			{ID: 123, WalletID: 123},
			{ID: 456, WalletID: 123},
		}, entries)
	})

	t.Run("Wallet doesn't belongs to the user", func(t *testing.T) {
		wallet := &domain.Wallet{UserID: 456}
		mockWalletRepo.EXPECT().Get(ctx, args.WalletID).Return(wallet, nil)

		_, err := u.List(ctx, args)
		assert.EqualError(t, err, "wallet doesn't belong to the authenticated user")
	})
}
