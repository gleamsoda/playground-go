package wallet

import (
	"context"
)

type Usecase interface {
	CreateAccount(ctx context.Context, args *CreateAccountParams) (*Account, error)
	GetAccount(ctx context.Context, args *GetAccountsParams) (*Account, error)
	ListAccounts(ctx context.Context, args *ListAccountsParams) ([]Account, error)
	CreateTransfer(ctx context.Context, args *CreateTransferParams) (*Transfer, error)
	CreateUser(ctx context.Context, args *CreateUserParams) (*User, error)
	LoginUser(ctx context.Context, args *LoginUserParams) (*LoginUserOutputParams, error)
	UpdateUser(ctx context.Context, args *UpdateUserParams) (*User, error)
	RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenOutputParams, error)
	SendVerifyEmail(ctx context.Context, args *SendVerifyEmailPayload) (*VerifyEmail, error)
	VerifyEmail(ctx context.Context, args *VerifyEmailParams) (*User, error)
}
