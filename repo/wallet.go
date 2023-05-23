package repo

import (
	"context"
	"database/sql"

	"github.com/gleamsoda/go-playground/domain"
	"github.com/gleamsoda/go-playground/repo/internal/sqlc"
)

type WalletRepository struct {
	q  sqlc.Querier
	db *sql.DB
}

var _ domain.WalletRepository = (*WalletRepository)(nil)

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{
		q:  sqlc.New(db),
		db: db,
	}
}

func (r *WalletRepository) Create(ctx context.Context, arg domain.CreateWalletParams) (*domain.Wallet, error) {
	id, err := r.q.CreateWallet(ctx, sqlc.CreateWalletParams{
		UserID:   arg.UserID,
		Balance:  arg.Balance,
		Currency: arg.Currency,
	})
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

func (r *WalletRepository) Get(ctx context.Context, id int64) (*domain.Wallet, error) {
	wallet, err := r.q.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Wallet{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		Currency:  wallet.Currency,
		CreatedAt: wallet.CreatedAt,
	}, nil
}

func (r *WalletRepository) List(ctx context.Context, arg domain.ListWalletsParams) ([]domain.Wallet, error) {
	wallets, err := r.q.ListWallets(ctx, sqlc.ListWalletsParams{
		UserID: arg.UserID,
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var result []domain.Wallet
	for _, wallet := range wallets {
		result = append(result, domain.Wallet{
			ID:        wallet.ID,
			UserID:    wallet.UserID,
			Balance:   wallet.Balance,
			Currency:  wallet.Currency,
			CreatedAt: wallet.CreatedAt,
		})
	}

	return result, nil
}

func (r *WalletRepository) Delete(ctx context.Context, id int64) error {
	return r.q.DeleteWallet(ctx, id)
}
