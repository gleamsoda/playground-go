package resolver

import (
	"time"

	"github.com/samber/do"

	"playground/internal/app"
	"playground/internal/app/usecase"
	"playground/internal/pkg/token"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	createUser    app.CreateUserUsecase
	loginUser     app.LoginUserUsecase
	createAccount app.CreateAccountUsecase
	listAccounts  app.ListAccountsUsecase
}

func NewResolver(i *do.Injector) (*Resolver, error) {
	r := do.MustInvoke[app.Repository](i)
	d := do.MustInvoke[app.Dispatcher](i)
	tm := do.MustInvoke[token.Manager](i)
	accessTokenDuration := do.MustInvokeNamed[time.Duration](i, "AccessTokenDuration")
	refreshTokenDuration := do.MustInvokeNamed[time.Duration](i, "RefreshTokenDuration")

	return &Resolver{
		createUser:    usecase.NewCreateUser(r, d),
		loginUser:     usecase.NewLoginUser(r, tm, accessTokenDuration, refreshTokenDuration),
		createAccount: usecase.NewCreateAccount(r),
		listAccounts:  usecase.NewListAccounts(r),
	}, nil
}
