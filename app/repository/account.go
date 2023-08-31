package repository

import (
	"context"
	"database/sql"

	"playground/app"
	"playground/app/repository/sqlc/gen"
)

type AccountRepository struct {
	db *sql.DB
	q  *gen.Queries
}

func NewAccountRepository(db *sql.DB) app.AccountRepository {
	return &AccountRepository{
		db: db,
		q:  gen.New(db),
	}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, args *app.Account) (*app.Account, error) {
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

func (r *AccountRepository) GetAccount(ctx context.Context, id int64) (*app.Account, error) {
	a, err := r.q.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	return &app.Account{
		ID:        a.ID,
		Owner:     a.Owner,
		Balance:   a.Balance,
		Currency:  a.Currency,
		CreatedAt: a.CreatedAt,
	}, nil
}

func (r *AccountRepository) ListAccounts(ctx context.Context, args *app.ListAccountsParams) ([]app.Account, error) {
	as, err := r.q.ListAccounts(ctx, &gen.ListAccountsParams{
		Owner:  args.Owner,
		Limit:  args.Limit,
		Offset: args.Offset,
	})
	if err != nil {
		return nil, err
	}

	result := []app.Account{}
	for _, a := range as {
		result = append(result, app.Account{
			ID:        a.ID,
			Owner:     a.Owner,
			Balance:   a.Balance,
			Currency:  a.Currency,
			CreatedAt: a.CreatedAt,
		})
	}

	return result, nil
}