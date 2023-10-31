package app

import (
	"context"
	"time"
)

type Transfer struct {
	ID            int64 `json:"id"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	// must be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func NewTransfer(fromAccountID, toAccountID, amount int64) *Transfer {
	return &Transfer{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
	}
}

type CreateTransferUsecase interface {
	Execute(ctx context.Context, args *CreateTransferParams) (*Transfer, error)
}

type CreateTransferParams struct {
	RequestUsername string `json:"-"`
	FromAccountID   int64  `json:"from_account_id"`
	ToAccountID     int64  `json:"to_account_id"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
}

type ListTransfersParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `json:"offset"`
}
