package repository

import (
	"context"
	"database/sql"

	"playground/internal/app"
	"playground/internal/app/repository/sqlc/gen"
)

type Transfer struct {
	exec Executor
	q    gen.Querier
}

func NewTransfer(e Executor) *Transfer {
	return &Transfer{
		exec: e,
		q:    gen.New(e),
	}
}

var _ app.TransferRepository = (*Transfer)(nil)

func (r *Transfer) Create(ctx context.Context, args *app.Transfer) (*app.Transfer, error) {
	var t *gen.Transfer

	if err := runTx(ctx, r.exec, func(cctx context.Context, tx *sql.Tx) error {
		txr := NewTransfer(tx)
		id, err := txr.q.CreateTransfer(cctx, &gen.CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		_, err = txr.q.CreateEntry(cctx, &gen.CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		_, err = txr.q.CreateEntry(cctx, &gen.CreateEntryParams{
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
			if err := txr.q.AddAccountBalance(cctx, b); err != nil {
				return err
			}
		}

		if t, err = txr.q.GetTransfer(cctx, id); err != nil {
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

func (r *Transfer) Get(ctx context.Context, id int64) (*app.Transfer, error) {
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

func (r *Transfer) List(ctx context.Context, args *app.ListTransfersParams) ([]app.Transfer, error) {
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

// func (r *Transfer) CreateEntry(ctx context.Context, args *app.Entry) (*app.Entry, error) {
// 	id, err := r.q.CreateEntry(ctx, &gen.CreateEntryParams{
// 		AccountID: args.AccountID,
// 		Amount:    args.Amount,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return r.GetEntry(ctx, id)
// }

// func (r *Transfer) GetEntry(ctx context.Context, id int64) (*app.Entry, error) {
// 	e, err := r.q.GetEntry(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &app.Entry{
// 		ID:        e.ID,
// 		AccountID: e.AccountID,
// 		Amount:    e.Amount,
// 		CreatedAt: e.CreatedAt,
// 	}, nil
// }
