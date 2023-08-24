package repository

import (
	"context"
	"database/sql"

	"playground/app"
	"playground/app/repository/gen"
)

type WalletRepository struct {
	q *gen.Queries
}

var _ app.WalletRepository = (*WalletRepository)(nil)

func NewWalletRepository(db *sql.DB) app.WalletRepository {
	return &WalletRepository{
		q: gen.New(db),
	}
}

func (r *WalletRepository) WithCtx(ctx context.Context) app.WalletRepository {
	if tx, ok := ctx.Value(TransactionKey).(*sql.Tx); ok {
		r.q.WithTx(tx)
	}
	return r
}

func (r *WalletRepository) Create(ctx context.Context, arg *app.Wallet) (*app.Wallet, error) {
	id, err := r.q.CreateWallet(ctx, gen.CreateWalletParams{
		UserID:   arg.UserID,
		Balance:  arg.Balance,
		Currency: arg.Currency,
	})
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

func (r *WalletRepository) Get(ctx context.Context, id int64) (*app.Wallet, error) {
	wallet, err := r.q.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	return &app.Wallet{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		Currency:  wallet.Currency,
		CreatedAt: wallet.CreatedAt,
	}, nil
}

func (r *WalletRepository) List(ctx context.Context, arg app.ListWalletsInputParams) ([]app.Wallet, error) {
	wallets, err := r.q.ListWallets(ctx, gen.ListWalletsParams{
		UserID: arg.UserID,
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var result []app.Wallet
	for _, wallet := range wallets {
		result = append(result, app.Wallet{
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

func (r *WalletRepository) AddWalletBalance(ctx context.Context, arg app.AddWalletBalanceParams) error {
	return r.q.AddWalletBalance(ctx, gen.AddWalletBalanceParams{
		ID:     arg.ID,
		Amount: arg.Amount,
	})
}
