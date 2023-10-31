package repository

import (
	"context"

	"playground/internal/app"
	"playground/internal/app/repository/sqlc/gen"
)

func (r *Repository) CreateTransfer(ctx context.Context, args *app.Transfer) (*app.Transfer, error) {
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

func (r *Repository) GetTransfer(ctx context.Context, id int64) (*app.Transfer, error) {
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

func (r *Repository) ListTransfers(ctx context.Context, args *app.ListTransfersParams) ([]app.Transfer, error) {
	ts, err := r.q.ListTransfers(ctx, &gen.ListTransfersParams{
		FromAccountID: args.FromAccountID,
		ToAccountID:   args.ToAccountID,
		Limit:         args.Limit,
		Offset:        args.Offset,
	})
	if err != nil {
		return nil, err
	}

	result := []app.Transfer{}
	for _, t := range ts {
		result = append(result, app.Transfer{
			ID:            t.ID,
			FromAccountID: t.FromAccountID,
			ToAccountID:   t.ToAccountID,
			Amount:        t.Amount,
			CreatedAt:     t.CreatedAt,
		})
	}

	return result, nil
}

func (r *Repository) CreateEntry(ctx context.Context, args *app.Entry) (*app.Entry, error) {
	id, err := r.q.CreateEntry(ctx, &gen.CreateEntryParams{
		AccountID: args.AccountID,
		Amount:    args.Amount,
	})
	if err != nil {
		return nil, err
	}

	return r.GetEntry(ctx, id)
}

func (r *Repository) GetEntry(ctx context.Context, id int64) (*app.Entry, error) {
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
