package app

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateAccount(ctx context.Context, args *Account) (*Account, error)
	GetAccount(ctx context.Context, id int64) (*Account, error)
	ListAccounts(ctx context.Context, args *ListAccountsParams) ([]Account, error)
	CreateTransfer(ctx context.Context, args *Transfer) (*Transfer, error)
	CreateUser(ctx context.Context, args *User) (*User, error)
	GetUser(ctx context.Context, username string) (*User, error)
	CreateSession(ctx context.Context, args *Session) error
	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
}
