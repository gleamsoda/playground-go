package repository

import (
	"context"
	"database/sql"

	"playground/app"
	"playground/app/repository/gen"
)

type TransferRepository struct {
	q *gen.Queries
}

var _ app.TransferRepository = (*TransferRepository)(nil)

func NewTransferRepository(db *sql.DB) app.TransferRepository {
	return &TransferRepository{
		q: gen.New(db),
	}
}

func (r *TransferRepository) WithCtx(ctx context.Context) app.TransferRepository {
	if tx, ok := ctx.Value(TransactionKey).(*sql.Tx); ok {
		r.q.WithTx(tx)
	}
	return r
}

func (r *TransferRepository) Create(ctx context.Context, arg *app.Transfer) (*app.Transfer, error) {
	id, err := r.q.CreateTransfer(ctx, gen.CreateTransferParams{
		FromWalletID: arg.FromWalletID,
		ToWalletID:   arg.ToWalletID,
		Amount:       arg.Amount,
	})
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

func (r *TransferRepository) Get(ctx context.Context, id int64) (*app.Transfer, error) {
	transfer, err := r.q.GetTransfer(ctx, id)
	if err != nil {
		return nil, err
	}

	return &app.Transfer{
		ID:           transfer.ID,
		FromWalletID: transfer.FromWalletID,
		ToWalletID:   transfer.ToWalletID,
		Amount:       transfer.Amount,
		CreatedAt:    transfer.CreatedAt,
	}, nil

}
