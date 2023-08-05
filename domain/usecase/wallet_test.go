package usecase

import (
	"context"
	"testing"

	"playground/domain"
	mock_domain "playground/test/mock/domain"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWalletUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mock_domain.NewMockWalletRepository(ctrl)
	u := &WalletUsecase{
		r: mockWalletRepo,
	}
	ctx := context.Background()
	args := domain.CreateWalletInputParams{
		UserID:   123,
		Balance:  100,
		Currency: "USD",
	}
	wallet := &domain.Wallet{
		UserID:   123,
		Balance:  100,
		Currency: "USD",
	}

	mockWalletRepo.EXPECT().Create(ctx, &domain.Wallet{
		UserID:   args.UserID,
		Balance:  args.Balance,
		Currency: args.Currency,
	}).Return(wallet, nil)

	got, err := u.Create(ctx, args)
	assert.NoError(t, err)
	assert.Equal(t, &domain.Wallet{
		UserID:   123,
		Balance:  100,
		Currency: "USD",
	}, got)
}

func TestWalletUsecase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mock_domain.NewMockWalletRepository(ctrl)
	u := &WalletUsecase{
		r: mockWalletRepo,
	}
	ctx := context.Background()
	args := domain.GetOrDeleteWalletInputParams{
		ID:     123,
		UserID: 123,
	}

	t.Run("Wallet belongs to the user", func(t *testing.T) {
		wallet := &domain.Wallet{
			ID:       123,
			UserID:   123,
			Balance:  100,
			Currency: "USD",
		}

		mockWalletRepo.EXPECT().Get(ctx, args.ID).Return(wallet, nil)

		got, err := u.Get(ctx, args)
		assert.NoError(t, err)
		assert.Equal(t, &domain.Wallet{
			ID:       123,
			UserID:   123,
			Balance:  100,
			Currency: "USD",
		}, got)
	})

	t.Run("Wallet doesn't belongs to the user", func(t *testing.T) {
		wallet := &domain.Wallet{
			ID:       123,
			UserID:   456,
			Balance:  100,
			Currency: "USD",
		}

		mockWalletRepo.EXPECT().Get(ctx, args.ID).Return(wallet, nil)

		_, err := u.Get(ctx, args)
		assert.Error(t, err)
	})
}

func TestWalletUsecase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mock_domain.NewMockWalletRepository(ctrl)
	u := &WalletUsecase{
		r: mockWalletRepo,
	}
	ctx := context.Background()
	args := domain.ListWalletsInputParams{
		UserID: 123,
		Limit:  10,
		Offset: 0,
	}
	walletList := []domain.Wallet{
		{ID: 123, UserID: 123, Balance: 100, Currency: "USD"},
		{ID: 456, UserID: 123, Balance: 100, Currency: "USD"},
	}

	mockWalletRepo.EXPECT().List(ctx, args).Return(walletList, nil)

	got, err := u.List(ctx, args)
	assert.NoError(t, err)
	assert.Equal(t, []domain.Wallet{
		{ID: 123, UserID: 123, Balance: 100, Currency: "USD"},
		{ID: 456, UserID: 123, Balance: 100, Currency: "USD"},
	}, got)
}

func TestWalletUsecase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mock_domain.NewMockWalletRepository(ctrl)
	u := &WalletUsecase{
		r: mockWalletRepo,
	}
	ctx := context.Background()
	args := domain.GetOrDeleteWalletInputParams{
		ID:     123,
		UserID: 123,
	}

	t.Run("Wallet belongs to the user", func(t *testing.T) {
		wallet := &domain.Wallet{
			ID:       123,
			UserID:   123,
			Balance:  100,
			Currency: "USD",
		}

		mockWalletRepo.EXPECT().Get(ctx, args.ID).Return(wallet, nil)
		mockWalletRepo.EXPECT().Delete(ctx, wallet.ID).Return(nil)

		err := u.Delete(ctx, args)
		assert.NoError(t, err)
	})

	t.Run("Wallet doesn't belongs to the user", func(t *testing.T) {
		wallet := &domain.Wallet{
			ID:       123,
			UserID:   456,
			Balance:  100,
			Currency: "USD",
		}

		mockWalletRepo.EXPECT().Get(ctx, args.ID).Return(wallet, nil)

		err := u.Delete(ctx, args)
		assert.Error(t, err)
	})
}
