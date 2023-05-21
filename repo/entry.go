package repo

import (
	"context"
	"database/sql"

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

	return r.GetByID(ctx, id)
}

func (r *EntryRepository) GetByID(ctx context.Context, id int64) (*domain.Entry, error) {
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
