package domain

import (
	"context"
	"time"
)

type Transfer struct {
	ID           int64     `json:"id"`
	FromWalletID int64     `json:"from_wallet_id"`
	ToWalletID   int64     `json:"to_wallet_id"`
	Amount       int64     `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
}

type TransferUsecase interface {
	Create(ctx context.Context, arg CreateTransferInputParams) (*Transfer, error)
}

type TransferRepository interface {
	Create(ctx context.Context, arg *Transfer) (*Transfer, error)
	WithCtx(ctx context.Context) TransferRepository
}

func NewTransfer(fromWalletID, toWalletID, amount int64) *Transfer {
	return &Transfer{
		FromWalletID: fromWalletID,
		ToWalletID:   toWalletID,
		Amount:       amount,
	}
}

type CreateTransferInputParams struct {
	RequestUserID int64  `json:"-"`
	FromWalletID  int64  `json:"from_wallet_id"`
	ToWalletID    int64  `json:"to_wallet_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}
