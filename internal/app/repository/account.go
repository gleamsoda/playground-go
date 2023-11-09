package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/morikuni/failure"

	"playground/internal/app"
	"playground/internal/app/repository/sqlc/gen"
	"playground/internal/pkg/apperr"
)

type Account struct {
	q gen.Querier
}

func NewAccount(e Executor) *Account {
	return &Account{
		q: gen.New(e),
	}
}

var _ app.AccountRepository = (*Account)(nil)

func (r *Account) Create(ctx context.Context, args *app.Account) (*app.Account, error) {
	id, err := r.q.CreateAccount(ctx, &gen.CreateAccountParams{
		Owner:    args.Owner,
		Balance:  args.Balance,
		Currency: args.Currency,
	})
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

func (r *Account) Get(ctx context.Context, id int64) (*app.Account, error) {
	a, err := r.q.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.Translate(err, apperr.NotFound)
		}
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

func (r *Account) List(ctx context.Context, args *app.ListAccountsParams) ([]app.Account, error) {
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

func (r *Account) Update(ctx context.Context, args *app.Account) (*app.Account, error) {
	err := r.q.UpdateAccount(ctx, &gen.UpdateAccountParams{
		ID:      args.ID,
		Balance: args.Balance,
	})
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, args.ID)
}

func (r *Account) Delete(ctx context.Context, id int64) error {
	return r.q.DeleteAccount(ctx, id)
}
