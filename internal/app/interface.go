package app

import (
	"context"

	"github.com/google/uuid"
)

type ( // Repository is an interface for database access.
	Repository interface {
		Transaction(ctx context.Context, fn TransactionFunc) error
		CreateAccount(ctx context.Context, args *Account) (*Account, error)
		GetAccount(ctx context.Context, id int64) (*Account, error)
		ListAccounts(ctx context.Context, args *ListAccountsParams) ([]Account, error)
		UpdateAccount(ctx context.Context, args *Account) (*Account, error)
		DeleteAccount(ctx context.Context, id int64) error
		CreateTransfer(ctx context.Context, args *Transfer) (*Transfer, error)
		GetTransfer(ctx context.Context, id int64) (*Transfer, error)
		ListTransfers(ctx context.Context, args *ListTransfersParams) ([]Transfer, error)
		CreateUser(ctx context.Context, args *User) (*User, error)
		GetUser(ctx context.Context, username string) (*User, error)
		UpdateUser(ctx context.Context, args *User) (*User, error)
		CreateSession(ctx context.Context, args *Session) error
		GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
		CreateVerifyEmail(ctx context.Context, args *VerifyEmail) (*VerifyEmail, error)
		UpdateVerifyEmail(ctx context.Context, args *VerifyEmail) (*VerifyEmail, error)
		UpdateUserEmailVerified(ctx context.Context, args *VerifyEmailParams) (*User, *VerifyEmail, error)
	}
	TransactionFunc     func(context.Context, Repository) error
	ListTransfersParams struct {
		FromAccountID int64 `json:"from_account_id"`
		ToAccountID   int64 `json:"to_account_id"`
		Limit         int32 `json:"limit"`
		Offset        int32 `json:"offset"`
	}
)

const SendVerifyEmailQueue = "task:send_verify_email"

type ( // Dispatcher is an interface for sending messages to a queue.
	Dispatcher interface {
		SendVerifyEmail(ctx context.Context, payload *SendVerifyEmailPayload) error
	}
	SendVerifyEmailPayload struct {
		Username string `json:"username" binding:"required,alphanum"`
	}
)
