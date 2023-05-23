package domain

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
	Create(ctx context.Context, arg CreateEntryParams) (*Entry, error)
	Get(ctx context.Context, id int64) (*Entry, error)
	List(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
}

type EntryRepository interface {
	Create(ctx context.Context, arg CreateEntryParams) (*Entry, error)
	Get(ctx context.Context, id int64) (*Entry, error)
	List(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
}

type CreateEntryParams struct {
	WalletID int64 `json:"wallet_id"`
	Amount   int64 `json:"amount"`
}

type ListEntriesParams struct {
	WalletID int64 `json:"wallet_id"`
	Limit    int32 `json:"limit" form:"limit"`
	Offset   int32 `json:"offset" form:"offset"`
}
