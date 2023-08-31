package app

import "context"

type Usecase interface {
	CreateAccount(ctx context.Context, args *CreateAccountParams) (*Account, error)
	GetAccount(ctx context.Context, id int64) (*Account, error)
	ListAccounts(ctx context.Context, args *ListAccountsParams) ([]Account, error)
	CreateTransfer(ctx context.Context, args *CreateTransferParams) (*Transfer, error)
	CreateUser(ctx context.Context, args *CreateUserParams) (*User, error)
	LoginUser(ctx context.Context, args *LoginUserParams) (*LoginUserOutputParams, error)
	RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenOutputParams, error)
}
