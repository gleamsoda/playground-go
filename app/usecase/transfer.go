package usecase

import (
	"context"
	"errors"
	"fmt"

	"playground/app"
)

type TransferUsecase struct {
	transferRepo app.TransferRepository
	entryRepo    app.EntryRepository
	walletRepo   app.WalletRepository
	txManager    app.TransactionManager
}

var _ app.TransferUsecase = (*TransferUsecase)(nil)

func NewTransferUsecase(transferRepo app.TransferRepository, entryRepo app.EntryRepository, walletRepo app.WalletRepository, txManager app.TransactionManager) app.TransferUsecase {
	return &TransferUsecase{
		transferRepo: transferRepo,
		entryRepo:    entryRepo,
		walletRepo:   walletRepo,
		txManager:    txManager,
	}
}

func (u *TransferUsecase) Create(ctx context.Context, arg app.CreateTransferInputParams) (*app.Transfer, error) {
	fromWallet, err := u.validWallet(ctx, arg.FromWalletID, arg.Currency)
	if err != nil {
		return nil, err
	}
	if fromWallet.UserID != arg.RequestUserID {
		err := errors.New("from wallet doesn't belong to the authenticated user")
		return nil, err
	}
	_, err = u.validWallet(ctx, arg.ToWalletID, arg.Currency)
	if err != nil {
		return nil, err
	}

	var t *app.Transfer

	if err := u.txManager.Do(ctx, func(innerCtx context.Context) error {
		var err error
		tRepo := u.transferRepo.WithCtx(innerCtx)
		eRepo := u.entryRepo.WithCtx(innerCtx)

		t, err = tRepo.Create(innerCtx, app.NewTransfer(arg.FromWalletID, arg.ToWalletID, arg.Amount))
		if err != nil {
			return err
		}

		if _, err = eRepo.Create(innerCtx, app.NewEntry(t.FromWalletID, -t.Amount)); err != nil {
			return err
		}
		if _, err := eRepo.Create(innerCtx, app.NewEntry(t.ToWalletID, t.Amount)); err != nil {
			return err
		}

		// avoid deadlock
		if arg.FromWalletID < arg.ToWalletID {
			err = u.addMoney(innerCtx, t.FromWalletID, -t.Amount, t.ToWalletID, t.Amount)
		} else {
			err = u.addMoney(innerCtx, t.ToWalletID, t.Amount, t.FromWalletID, -t.Amount)
		}
		return err
	}); err != nil {
		return nil, err
	}

	return t, err
}

func (u *TransferUsecase) addMoney(ctx context.Context, walletID1 int64, amount1 int64, walletID2 int64, amount2 int64) error {
	wRepo := u.walletRepo.WithCtx(ctx)
	if err := wRepo.AddWalletBalance(ctx, app.AddWalletBalanceParams{
		ID:     walletID1,
		Amount: amount1,
	}); err != nil {
		return err
	}
	return wRepo.AddWalletBalance(ctx, app.AddWalletBalanceParams{
		ID:     walletID2,
		Amount: amount2,
	})
}

func (u *TransferUsecase) validWallet(ctx context.Context, walletD int64, currency string) (*app.Wallet, error) {
	w, err := u.walletRepo.Get(ctx, walletD)
	if err != nil {
		return nil, err
	}

	if w.Currency != currency {
		err := fmt.Errorf("w [%d] currency mismatch: %s vs %s", w.ID, w.Currency, currency)
		return nil, err
	}

	return w, nil
}
