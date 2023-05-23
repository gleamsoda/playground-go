package domain

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
	Create(ctx context.Context, arg CreateWalletParams) (*Wallet, error)
	Get(ctx context.Context, id int64) (*Wallet, error)
	List(ctx context.Context, arg ListWalletsParams) ([]Wallet, error)
	Delete(ctx context.Context, id int64) error
}

type WalletRepository interface {
	Create(ctx context.Context, arg CreateWalletParams) (*Wallet, error)
	Get(ctx context.Context, id int64) (*Wallet, error)
	List(ctx context.Context, arg ListWalletsParams) ([]Wallet, error)
	Delete(ctx context.Context, id int64) error
}

type CreateWalletParams struct {
	UserID   int64  `json:"user_id"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type ListWalletsParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit" form:"limit"`
	Offset int32 `json:"offset" form:"offset"`
}
