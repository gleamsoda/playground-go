package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"playground/app"
	mock_app "playground/test/mock/app"
)

func TestTransferUsecase_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_app.NewMockTransferRepository(ctrl)
	mockEntryRepo := mock_app.NewMockEntryRepository(ctrl)
	mockWalletRepo := mock_app.NewMockWalletRepository(ctrl)
	u := &TransferUsecase{
		transferRepo: mockTransferRepo,
		entryRepo:    mockEntryRepo,
		walletRepo:   mockWalletRepo,
		txManager:    &app.PassThroughTxManager{},
	}
	args := app.CreateTransferInputParams{
		RequestUserID: 123,
		FromWalletID:  123,
		ToWalletID:    456,
		Amount:        100,
		Currency:      "USD",
	}

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		fromWallet := &app.Wallet{UserID: 123, Currency: "USD"}
		toWallet := &app.Wallet{UserID: 456, Currency: "USD"}
		tr := &app.Transfer{
			FromWalletID: 123,
			ToWalletID:   456,
			Amount:       100,
		}

		mockWalletRepo.EXPECT().Get(ctx, args.FromWalletID).Return(fromWallet, nil)
		mockWalletRepo.EXPECT().Get(ctx, args.ToWalletID).Return(toWallet, nil)
		mockTransferRepo.EXPECT().WithCtx(ctx).Return(mockTransferRepo)
		mockEntryRepo.EXPECT().WithCtx(ctx).Return(mockEntryRepo)
		mockTransferRepo.EXPECT().Create(gomock.Any(), tr).Return(tr, nil)
		mockEntryRepo.EXPECT().Create(gomock.Any(), &app.Entry{WalletID: tr.FromWalletID, Amount: -tr.Amount}).Return(&app.Entry{WalletID: tr.FromWalletID, Amount: -tr.Amount}, nil)
		mockEntryRepo.EXPECT().Create(gomock.Any(), &app.Entry{WalletID: tr.ToWalletID, Amount: tr.Amount}).Return(&app.Entry{WalletID: tr.ToWalletID, Amount: tr.Amount}, nil)
		mockWalletRepo.EXPECT().WithCtx(gomock.Any()).Return(mockWalletRepo)
		gomock.InOrder(
			mockWalletRepo.EXPECT().AddWalletBalance(gomock.Any(), app.AddWalletBalanceParams{
				ID:     tr.FromWalletID,
				Amount: -tr.Amount,
			}).Return(nil),
			mockWalletRepo.EXPECT().AddWalletBalance(gomock.Any(), app.AddWalletBalanceParams{
				ID:     tr.ToWalletID,
				Amount: tr.Amount,
			}).Return(nil),
		)

		transfer, err := u.Create(context.Background(), args)

		assert.NoError(t, err)
		assert.Equal(t, tr, transfer)
	})

	t.Run("Currency Mismatch fromWallet and args", func(t *testing.T) {
		fromWallet := &app.Wallet{UserID: 123, Currency: "EUR"}

		mockWalletRepo.EXPECT().Get(ctx, args.FromWalletID).Return(fromWallet, nil)

		_, err := u.Create(context.Background(), args)
		assert.Error(t, err)
	})

	t.Run("UserID Mismatch fromWallet and args", func(t *testing.T) {
		fromWallet := &app.Wallet{UserID: 456, Currency: "USD"}

		mockWalletRepo.EXPECT().Get(ctx, args.FromWalletID).Return(fromWallet, nil)

		_, err := u.Create(context.Background(), args)
		assert.Error(t, err)
	})

	t.Run("Currency Mismatch fromWallet and toWallet", func(t *testing.T) {
		fromWallet := &app.Wallet{UserID: 123, Currency: "USD"}
		toWallet := &app.Wallet{UserID: 456, Currency: "EUR"}

		mockWalletRepo.EXPECT().Get(ctx, args.FromWalletID).Return(fromWallet, nil)
		mockWalletRepo.EXPECT().Get(ctx, args.ToWalletID).Return(toWallet, nil)

		_, err := u.Create(context.Background(), args)
		assert.Error(t, err)
	})

	t.Run("Wrong addMoney order", func(t *testing.T) {
		args := app.CreateTransferInputParams{
			RequestUserID: 456,
			FromWalletID:  456,
			ToWalletID:    123,
			Amount:        100,
			Currency:      "USD",
		}
		fromWallet := &app.Wallet{UserID: 456, Currency: "USD"}
		toWallet := &app.Wallet{UserID: 123, Currency: "USD"}
		tr := &app.Transfer{
			FromWalletID: 456,
			ToWalletID:   123,
			Amount:       100,
		}

		mockWalletRepo.EXPECT().Get(ctx, args.FromWalletID).Return(fromWallet, nil)
		mockWalletRepo.EXPECT().Get(ctx, args.ToWalletID).Return(toWallet, nil)
		mockTransferRepo.EXPECT().WithCtx(ctx).Return(mockTransferRepo)
		mockEntryRepo.EXPECT().WithCtx(ctx).Return(mockEntryRepo)
		mockTransferRepo.EXPECT().Create(gomock.Any(), tr).Return(tr, nil)
		mockEntryRepo.EXPECT().Create(gomock.Any(), &app.Entry{WalletID: tr.FromWalletID, Amount: -tr.Amount}).Return(&app.Entry{WalletID: tr.FromWalletID, Amount: -tr.Amount}, nil)
		mockEntryRepo.EXPECT().Create(gomock.Any(), &app.Entry{WalletID: tr.ToWalletID, Amount: tr.Amount}).Return(&app.Entry{WalletID: tr.ToWalletID, Amount: tr.Amount}, nil)
		mockWalletRepo.EXPECT().WithCtx(gomock.Any()).Return(mockWalletRepo)
		gomock.InOrder(
			mockWalletRepo.EXPECT().AddWalletBalance(gomock.Any(), app.AddWalletBalanceParams{
				ID:     tr.ToWalletID,
				Amount: tr.Amount,
			}).Return(nil),
			mockWalletRepo.EXPECT().AddWalletBalance(gomock.Any(), app.AddWalletBalanceParams{
				ID:     tr.FromWalletID,
				Amount: -tr.Amount,
			}).Return(nil),
		)

		transfer, err := u.Create(context.Background(), args)

		assert.NoError(t, err)
		assert.Equal(t, tr, transfer)
	})
}
