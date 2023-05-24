package repo

import (
	"context"
	"database/sql"

	"playground/domain"
	"playground/repo/sqlc/internal/boundary"
)

type TransferRepository struct {
	q *boundary.Queries
}

var _ domain.TransferRepository = (*TransferRepository)(nil)

func NewTransferRepository(db *sql.DB) domain.TransferRepository {
	return &TransferRepository{
		q: boundary.New(db),
	}
}

func (r *TransferRepository) WithCtx(ctx context.Context) domain.TransferRepository {
	if tx, ok := ctx.Value(TransactionKey).(*sql.Tx); ok {
		r.q.WithTx(tx)
	}
	return r
}

func (r *TransferRepository) Create(ctx context.Context, arg *domain.Transfer) (*domain.Transfer, error) {
	id, err := r.q.CreateTransfer(ctx, boundary.CreateTransferParams{
		FromWalletID: arg.FromWalletID,
		ToWalletID:   arg.ToWalletID,
		Amount:       arg.Amount,
	})
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

func (r *TransferRepository) Get(ctx context.Context, id int64) (*domain.Transfer, error) {
	transfer, err := r.q.GetTransfer(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Transfer{
		ID:           transfer.ID,
		FromWalletID: transfer.FromWalletID,
		ToWalletID:   transfer.ToWalletID,
		Amount:       transfer.Amount,
		CreatedAt:    transfer.CreatedAt,
	}, nil

}
