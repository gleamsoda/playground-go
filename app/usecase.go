package app

import "context"

type AccountUsecase interface {
	CreateAccount(ctx context.Context, args *CreateAccountParams) (*Account, error)
	GetAccount(ctx context.Context, id int64) (*Account, error)
	ListAccounts(ctx context.Context, args *ListAccountsParams) ([]Account, error)
}

type TransferUsecase interface {
	CreateTransfer(ctx context.Context, args *CreateTransferParams) (*Transfer, error)
}

type UserUsecase interface {
	CreateUser(ctx context.Context, args *CreateUserParams) (*User, error)
	Login(ctx context.Context, args *LoginUserParams) (*LoginUserOutputParams, error)
	RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenOutputParams, error)
}
