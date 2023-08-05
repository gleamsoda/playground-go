package usecase

import (
	"context"
	"testing"

	"playground/domain"
	mock_domain "playground/test/mock/domain"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTransferUsecase_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_domain.NewMockTransferRepository(ctrl)
	mockEntryRepo := mock_domain.NewMockEntryRepository(ctrl)
	mockWalletRepo := mock_domain.NewMockWalletRepository(ctrl)
	u := &TransferUsecase{
		transferRepo: mockTransferRepo,
		entryRepo:    mockEntryRepo,
		walletRepo:   mockWalletRepo,
		txManager:    &domain.PassThroughTxManager{},
	}
	args := domain.CreateTransferInputParams{
		RequestUserID: 123,
		FromWalletID:  123,
		ToWalletID:    456,
		Amount:        100,
		Currency:      "USD",
	}

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		fromWallet := &domain.Wallet{UserID: 123, Currency: "USD"}
		toWallet := &domain.Wallet{UserID: 456, Currency: "USD"}
		tr := &domain.Transfer{
			FromWalletID: 123,
			ToWalletID:   456,
			Amount:       100,
		}

		mockWalletRepo.EXPECT().Get(ctx, args.FromWalletID).Return(fromWallet, nil)
		mockWalletRepo.EXPECT().Get(ctx, args.ToWalletID).Return(toWallet, nil)
		mockTransferRepo.EXPECT().WithCtx(ctx).Return(mockTransferRepo)
		mockEntryRepo.EXPECT().WithCtx(ctx).Return(mockEntryRepo)
		mockTransferRepo.EXPECT().Create(gomock.Any(), tr).Return(tr, nil)
		mockEntryRepo.EXPECT().Create(gomock.Any(), &domain.Entry{WalletID: tr.FromWalletID, Amount: -tr.Amount}).Return(&domain.Entry{WalletID: tr.FromWalletID, Amount: -tr.Amount}, nil)
		mockEntryRepo.EXPECT().Create(gomock.Any(), &domain.Entry{WalletID: tr.ToWalletID, Amount: tr.Amount}).Return(&domain.Entry{WalletID: tr.ToWalletID, Amount: tr.Amount}, nil)
		mockWalletRepo.EXPECT().WithCtx(gomock.Any()).Return(mockWalletRepo)
		gomock.InOrder(
			mockWalletRepo.EXPECT().AddWalletBalance(gomock.Any(), domain.AddWalletBalanceParams{
				ID:     tr.FromWalletID,
				Amount: -tr.Amount,
			}).Return(nil),
			mockWalletRepo.EXPECT().AddWalletBalance(gomock.Any(), domain.AddWalletBalanceParams{
				ID:     tr.ToWalletID,
				Amount: tr.Amount,
			}).Return(nil),
		)

		transfer, err := u.Create(context.Background(), args)

		assert.NoError(t, err)
		assert.Equal(t, tr, transfer)
	})

	t.Run("Currency Mismatch fromWallet and args", func(t *testing.T) {
		fromWallet := &domain.Wallet{UserID: 123, Currency: "EUR"}

		mockWalletRepo.EXPECT().Get(ctx, args.FromWalletID).Return(fromWallet, nil)

		_, err := u.Create(context.Background(), args)
		assert.Error(t, err)
	})

	t.Run("UserID Mismatch fromWallet and args", func(t *testing.T) {
		fromWallet := &domain.Wallet{UserID: 456, Currency: "USD"}

		mockWalletRepo.EXPECT().Get(ctx, args.FromWalletID).Return(fromWallet, nil)

		_, err := u.Create(context.Background(), args)
		assert.Error(t, err)
	})

	t.Run("Currency Mismatch fromWallet and toWallet", func(t *testing.T) {
		fromWallet := &domain.Wallet{UserID: 123, Currency: "USD"}
		toWallet := &domain.Wallet{UserID: 456, Currency: "EUR"}

		mockWalletRepo.EXPECT().Get(ctx, args.FromWalletID).Return(fromWallet, nil)
		mockWalletRepo.EXPECT().Get(ctx, args.ToWalletID).Return(toWallet, nil)

		_, err := u.Create(context.Background(), args)
		assert.Error(t, err)
	})

	t.Run("Wrong addMoney order", func(t *testing.T) {
		args := domain.CreateTransferInputParams{
			RequestUserID: 456,
			FromWalletID:  456,
			ToWalletID:    123,
			Amount:        100,
			Currency:      "USD",
		}
		fromWallet := &domain.Wallet{UserID: 456, Currency: "USD"}
		toWallet := &domain.Wallet{UserID: 123, Currency: "USD"}
		tr := &domain.Transfer{
			FromWalletID: 456,
			ToWalletID:   123,
			Amount:       100,
		}

		mockWalletRepo.EXPECT().Get(ctx, args.FromWalletID).Return(fromWallet, nil)
		mockWalletRepo.EXPECT().Get(ctx, args.ToWalletID).Return(toWallet, nil)
		mockTransferRepo.EXPECT().WithCtx(ctx).Return(mockTransferRepo)
		mockEntryRepo.EXPECT().WithCtx(ctx).Return(mockEntryRepo)
		mockTransferRepo.EXPECT().Create(gomock.Any(), tr).Return(tr, nil)
		mockEntryRepo.EXPECT().Create(gomock.Any(), &domain.Entry{WalletID: tr.FromWalletID, Amount: -tr.Amount}).Return(&domain.Entry{WalletID: tr.FromWalletID, Amount: -tr.Amount}, nil)
		mockEntryRepo.EXPECT().Create(gomock.Any(), &domain.Entry{WalletID: tr.ToWalletID, Amount: tr.Amount}).Return(&domain.Entry{WalletID: tr.ToWalletID, Amount: tr.Amount}, nil)
		mockWalletRepo.EXPECT().WithCtx(gomock.Any()).Return(mockWalletRepo)
		gomock.InOrder(
			mockWalletRepo.EXPECT().AddWalletBalance(gomock.Any(), domain.AddWalletBalanceParams{
				ID:     tr.ToWalletID,
				Amount: tr.Amount,
			}).Return(nil),
			mockWalletRepo.EXPECT().AddWalletBalance(gomock.Any(), domain.AddWalletBalanceParams{
				ID:     tr.FromWalletID,
				Amount: -tr.Amount,
			}).Return(nil),
		)

		transfer, err := u.Create(context.Background(), args)

		assert.NoError(t, err)
		assert.Equal(t, tr, transfer)
	})
}
