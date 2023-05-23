package repo

import (
	"context"
	"database/sql"

	"github.com/gleamsoda/go-playground/domain"
	"github.com/gleamsoda/go-playground/repo/internal/sqlc"
)

type TransferRepository struct {
	q  sqlc.Querier
	db *sql.DB
}

var _ domain.TransferRepository = (*TransferRepository)(nil)

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{
		q:  sqlc.New(db),
		db: db,
	}
}

func (r *TransferRepository) Create(ctx context.Context, arg domain.CreateTransferParams) (*domain.Transfer, error) {
	id, err := r.q.CreateTransfer(ctx, sqlc.CreateTransferParams{
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
