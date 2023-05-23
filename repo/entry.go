package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gleamsoda/go-playground/domain"
	"github.com/gleamsoda/go-playground/repo/internal/sqlc"
)

type EntryRepository struct {
	q  sqlc.Querier
	db *sql.DB
}

var _ domain.EntryRepository = (*EntryRepository)(nil)

func NewEntryRepository(db *sql.DB) *EntryRepository {
	return &EntryRepository{
		q:  sqlc.New(db),
		db: db,
	}
}

func (r *EntryRepository) Create(ctx context.Context, arg domain.CreateEntryParams) (*domain.Entry, error) {
	id, err := r.q.CreateEntry(ctx, sqlc.CreateEntryParams{
		WalletID: arg.WalletID,
		Amount:   arg.Amount,
	})
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

func (r *EntryRepository) Get(ctx context.Context, id int64) (*domain.Entry, error) {
	entry, err := r.q.GetEntry(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Entry{
		ID:        entry.ID,
		WalletID:  entry.WalletID,
		Amount:    entry.Amount,
		CreatedAt: entry.CreatedAt,
	}, nil
}

func (r *EntryRepository) List(ctx context.Context, arg domain.ListEntriesParams) ([]domain.Entry, error) {
	fmt.Println("ListEntriesParams", arg)
	entries, err := r.q.ListEntries(ctx, sqlc.ListEntriesParams{
		WalletID: arg.WalletID,
		Limit:    arg.Limit,
		Offset:   arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var result []domain.Entry
	for _, entry := range entries {
		result = append(result, domain.Entry{
			ID:        entry.ID,
			WalletID:  entry.WalletID,
			Amount:    entry.Amount,
			CreatedAt: entry.CreatedAt,
		})
	}

	return result, nil
}
