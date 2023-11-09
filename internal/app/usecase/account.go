package usecase

import (
	"context"
	"errors"

	"github.com/morikuni/failure"

	"playground/internal/app"
	"playground/internal/pkg/apperr"
)

type (
	CreateAccount struct {
		r app.RepositoryManager
	}
	GetAccount struct {
		r app.RepositoryManager
	}
	ListAccounts struct {
		r app.RepositoryManager
	}
)

func NewCreateAccount(r app.RepositoryManager) *CreateAccount {
	return &CreateAccount{
		r: r,
	}
}

func (u *CreateAccount) Execute(ctx context.Context, args *app.CreateAccountParams) (*app.Account, error) {
	return u.r.Account().Create(ctx, app.NewAccount(args.Owner, args.Balance, args.Currency))
}

func NewGetAccount(r app.RepositoryManager) *GetAccount {
	return &GetAccount{
		r: r,
	}
}

func (u *GetAccount) Execute(ctx context.Context, args *app.GetAccountsParams) (*app.Account, error) {
	a, err := u.r.Account().Get(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if a.Owner != args.Owner {
		err := errors.New("account doesn't belong to the authenticated user")
		return nil, failure.Translate(err, apperr.Unauthenticated)
	}
	return a, nil
}

func NewListAccounts(r app.RepositoryManager) *ListAccounts {
	return &ListAccounts{
		r: r,
	}
}

func (u *ListAccounts) Execute(ctx context.Context, args *app.ListAccountsParams) ([]app.Account, error) {
	return u.r.Account().List(ctx, args)
}
