//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package app

import (
	"context"
)

type ( // Repository is an interface for database access.
	RepositoryManager interface {
		Transaction() Transaction
		Account() AccountRepository
		Transfer() TransferRepository
		User() UserRepository
	}
	Transaction interface {
		Run(ctx context.Context, fn TransactionFunc) error
	}
	TransactionFunc     func(context.Context, RepositoryManager) error
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
