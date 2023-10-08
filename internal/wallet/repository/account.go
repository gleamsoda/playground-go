package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/morikuni/failure"

	"playground/internal/pkg/apperr"
	"playground/internal/wallet"
	"playground/internal/wallet/repository/sqlc/gen"
)

func (r *Repository) CreateAccount(ctx context.Context, args *wallet.Account) (*wallet.Account, error) {
	id, err := r.q.CreateAccount(ctx, &gen.CreateAccountParams{
		Owner:    args.Owner,
		Balance:  args.Balance,
		Currency: args.Currency,
	})
	if err != nil {
		return nil, err
	}

	return r.GetAccount(ctx, id)
}

func (r *Repository) GetAccount(ctx context.Context, id int64) (*wallet.Account, error) {
	a, err := r.q.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.Translate(err, apperr.NotFound)
		}
		return nil, err
	}

	return &wallet.Account{
		ID:        a.ID,
		Owner:     a.Owner,
		Balance:   a.Balance,
		Currency:  a.Currency,
		CreatedAt: a.CreatedAt,
	}, nil
}

func (r *Repository) ListAccounts(ctx context.Context, args *wallet.ListAccountsParams) ([]wallet.Account, error) {
	as, err := r.q.ListAccounts(ctx, &gen.ListAccountsParams{
		Owner:  args.Owner,
		Limit:  args.Limit,
		Offset: args.Offset,
	})
	if err != nil {
		return nil, err
	}

	result := []wallet.Account{}
	for _, a := range as {
		result = append(result, wallet.Account{
			ID:        a.ID,
			Owner:     a.Owner,
			Balance:   a.Balance,
			Currency:  a.Currency,
			CreatedAt: a.CreatedAt,
		})
	}

	return result, nil
}

func (r *Repository) UpdateAccount(ctx context.Context, args *wallet.Account) (*wallet.Account, error) {
	err := r.q.UpdateAccount(ctx, &gen.UpdateAccountParams{
		ID:      args.ID,
		Balance: args.Balance,
	})
	if err != nil {
		return nil, err
	}
	return r.GetAccount(ctx, args.ID)
}

func (r *Repository) DeleteAccount(ctx context.Context, id int64) error {
	return r.q.DeleteAccount(ctx, id)
}
