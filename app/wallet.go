package app

import (
	"context"
	"time"
)

type Wallet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type WalletUsecase interface {
	Create(ctx context.Context, arg CreateWalletInputParams) (*Wallet, error)
	Get(ctx context.Context, arg GetOrDeleteWalletInputParams) (*Wallet, error)
	List(ctx context.Context, arg ListWalletsInputParams) ([]Wallet, error)
	Delete(ctx context.Context, arg GetOrDeleteWalletInputParams) error
}

type WalletRepository interface {
	Create(ctx context.Context, arg *Wallet) (*Wallet, error)
	Get(ctx context.Context, id int64) (*Wallet, error)
	List(ctx context.Context, arg ListWalletsInputParams) ([]Wallet, error)
	Delete(ctx context.Context, id int64) error
	AddWalletBalance(ctx context.Context, arg AddWalletBalanceParams) error
	WithCtx(ctx context.Context) WalletRepository
}

type GetOrDeleteWalletInputParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func NewWallet(userID int64, balance int64, currency string) *Wallet {
	return &Wallet{
		UserID:   userID,
		Balance:  balance,
		Currency: currency,
	}
}

type CreateWalletInputParams struct {
	UserID   int64  `json:"user_id"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency" binding:"required,currency"`
}

type ListWalletsInputParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit" form:"limit"`
	Offset int32 `json:"offset" form:"offset"`
}

type AddWalletBalanceParams struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"amount"`
}
