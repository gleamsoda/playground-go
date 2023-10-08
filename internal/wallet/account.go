package wallet

import (
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type GetAccountsParams struct {
	ID    int64  `json:"id"`
	Owner string `json:"owner"`
}

type ListAccountsParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func NewAccount(owner string, balance int64, currency string) *Account {
	return &Account{
		Owner:    owner,
		Balance:  balance,
		Currency: currency,
	}
}
