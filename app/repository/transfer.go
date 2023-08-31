package repository

import (
	"context"
	"database/sql"
	"fmt"

	"playground/app"
	"playground/app/repository/sqlc/gen"
)

type TransferRepository struct {
	db *sql.DB
	q  *gen.Queries
}

func NewTransferRepository(db *sql.DB) app.TransferRepository {
	return &TransferRepository{
		db: db,
		q:  gen.New(db),
	}
}

func (r *TransferRepository) CreateTransfer(ctx context.Context, args *app.Transfer) (*app.Transfer, error) {
	var t *gen.Transfer

	if err := r.tx(ctx, func(cctx context.Context, q *gen.Queries) error {
		id, err := q.CreateTransfer(cctx, &gen.CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateEntry(cctx, &gen.CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateEntry(cctx, &gen.CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		bs := []*gen.AddAccountBalanceParams{}
		if args.FromAccountID < args.ToAccountID {
			bs = append(bs, &gen.AddAccountBalanceParams{ID: args.FromAccountID, Amount: -args.Amount}, &gen.AddAccountBalanceParams{ID: args.ToAccountID, Amount: args.Amount})
		} else {
			bs = append(bs, &gen.AddAccountBalanceParams{ID: args.ToAccountID, Amount: args.Amount}, &gen.AddAccountBalanceParams{ID: args.FromAccountID, Amount: -args.Amount})
		}
		for _, b := range bs {
			if err := q.AddAccountBalance(cctx, b); err != nil {
				return err
			}
		}

		if t, err = q.GetTransfer(cctx, id); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &app.Transfer{
		ID:            t.ID,
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
		CreatedAt:     t.CreatedAt,
	}, nil
}

func (r *TransferRepository) GetTransfer(ctx context.Context, id int64) (*app.Transfer, error) {
	t, err := r.q.GetTransfer(ctx, id)
	if err != nil {
		return nil, err
	}

	return &app.Transfer{
		ID:            t.ID,
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
		CreatedAt:     t.CreatedAt,
	}, nil
}

func (r *TransferRepository) CreateEntry(ctx context.Context, args *app.Entry) (*app.Entry, error) {
	id, err := r.q.CreateEntry(ctx, &gen.CreateEntryParams{
		AccountID: args.AccountID,
		Amount:    args.Amount,
	})
	if err != nil {
		return nil, err
	}

	return r.GetEntry(ctx, id)
}

func (r *TransferRepository) GetEntry(ctx context.Context, id int64) (*app.Entry, error) {
	e, err := r.q.GetEntry(ctx, id)
	if err != nil {
		return nil, err
	}

	return &app.Entry{
		ID:        e.ID,
		AccountID: e.AccountID,
		Amount:    e.Amount,
		CreatedAt: e.CreatedAt,
	}, nil
}

func (r *TransferRepository) ListEntries(ctx context.Context, args app.ListEntriesParams) ([]app.Entry, error) {
	es, err := r.q.ListEntries(ctx, &gen.ListEntriesParams{
		AccountID: args.AccountID,
		Limit:     args.Limit,
		Offset:    args.Offset,
	})
	if err != nil {
		return nil, err
	}

	var result []app.Entry
	for _, e := range es {
		result = append(result, app.Entry{
			ID:        e.ID,
			AccountID: e.AccountID,
			Amount:    e.Amount,
			CreatedAt: e.CreatedAt,
		})
	}

	return result, nil
}

func (r *TransferRepository) tx(ctx context.Context, fn func(context.Context, *gen.Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := gen.New(tx)
	if err := fn(ctx, q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
