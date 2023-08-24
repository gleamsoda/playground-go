package app

import (
	"context"
	"time"
)

type Entry struct {
	ID        int64     `json:"id"`
	WalletID  int64     `json:"wallet_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type EntryUsecase interface {
	List(ctx context.Context, arg ListEntriesInputParams) ([]Entry, error)
}

type EntryRepository interface {
	Create(ctx context.Context, arg *Entry) (*Entry, error)
	Get(ctx context.Context, id int64) (*Entry, error)
	List(ctx context.Context, arg ListEntriesInputParams) ([]Entry, error)
	WithCtx(ctx context.Context) EntryRepository
}

func NewEntry(walletID, amount int64) *Entry {
	return &Entry{
		WalletID: walletID,
		Amount:   amount,
	}
}

type ListEntriesInputParams struct {
	RequestUserID int64 `json:"-"`
	WalletID      int64 `json:"wallet_id"`
	Limit         int32 `json:"limit" form:"limit"`
	Offset        int32 `json:"offset" form:"offset"`
}
