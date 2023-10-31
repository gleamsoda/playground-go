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
	createUserUsecase    app.CreateUserUsecase
	loginUserUsecase     app.LoginUserUsecase
	createAccountUsecase app.CreateAccountUsecase
	listAccountsUsecase  app.ListAccountsUsecase
}

func NewResolver(i *do.Injector) (*Resolver, error) {
	r := do.MustInvoke[app.Repository](i)
	d := do.MustInvoke[app.Dispatcher](i)
	tm := do.MustInvoke[token.Manager](i)
	accessTokenDuration := do.MustInvokeNamed[time.Duration](i, "AccessTokenDuration")
	refreshTokenDuration := do.MustInvokeNamed[time.Duration](i, "RefreshTokenDuration")

	return &Resolver{
		createUserUsecase:    usecase.NewCreateUserUsecase(r, d),
		loginUserUsecase:     usecase.NewLoginUserUsecase(r, tm, accessTokenDuration, refreshTokenDuration),
		createAccountUsecase: usecase.NewCreateAccountUsecase(r),
		listAccountsUsecase:  usecase.NewListAccountsUsecase(r),
	}, nil
}
