package repository

import (
	"context"
	"database/sql"

	"playground/app"
	"playground/app/repository/sqlc/gen"
)

type EntryRepository struct {
	q *gen.Queries
}

var _ app.EntryRepository = (*EntryRepository)(nil)

func NewEntryRepository(db *sql.DB) app.EntryRepository {
	return &EntryRepository{
		q: gen.New(db),
	}
}

func (r *EntryRepository) WithCtx(ctx context.Context) app.EntryRepository {
	if tx, ok := ctx.Value(TransactionKey).(*sql.Tx); ok {
		r.q.WithTx(tx)
	}
	return r
}

func (r *EntryRepository) Create(ctx context.Context, arg *app.Entry) (*app.Entry, error) {
	id, err := r.q.CreateEntry(ctx, gen.CreateEntryParams{
		WalletID: arg.WalletID,
		Amount:   arg.Amount,
	})
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

func (r *EntryRepository) Get(ctx context.Context, id int64) (*app.Entry, error) {
	entry, err := r.q.GetEntry(ctx, id)
	if err != nil {
		return nil, err
	}

	return &app.Entry{
		ID:        entry.ID,
		WalletID:  entry.WalletID,
		Amount:    entry.Amount,
		CreatedAt: entry.CreatedAt,
	}, nil
}

func (r *EntryRepository) List(ctx context.Context, arg app.ListEntriesInputParams) ([]app.Entry, error) {
	entries, err := r.q.ListEntries(ctx, gen.ListEntriesParams{
		WalletID: arg.WalletID,
		Limit:    arg.Limit,
		Offset:   arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var result []app.Entry
	for _, entry := range entries {
		result = append(result, app.Entry{
			ID:        entry.ID,
			WalletID:  entry.WalletID,
			Amount:    entry.Amount,
			CreatedAt: entry.CreatedAt,
		})
	}

	return result, nil
}
