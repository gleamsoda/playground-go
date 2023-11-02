//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package app

import (
	"context"
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(owner string, balance int64, currency string) *Account {
	return &Account{
		Owner:    owner,
		Balance:  balance,
		Currency: currency,
	}
}

type (
	CreateAccountUsecase interface {
		Execute(ctx context.Context, args *CreateAccountParams) (*Account, error)
	}
	CreateAccountParams struct {
		Owner    string `json:"owner"`
		Balance  int64  `json:"balance"`
		Currency string `json:"currency"`
	}
	GetAccountUsecase interface {
		Execute(ctx context.Context, args *GetAccountsParams) (*Account, error)
	}
	GetAccountsParams struct {
		ID    int64  `json:"id"`
		Owner string `json:"owner"`
	}
	ListAccountsUsecase interface {
		Execute(ctx context.Context, args *ListAccountsParams) ([]Account, error)
	}
	ListAccountsParams struct {
		Owner  string `json:"owner"`
		Limit  int32  `json:"limit"`
		Offset int32  `json:"offset"`
	}
)
