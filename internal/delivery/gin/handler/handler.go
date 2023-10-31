package handler

import (
	"time"

	"github.com/samber/do"

	"playground/internal/app"
	"playground/internal/app/usecase"
	"playground/internal/pkg/token"
)

type Handler struct {
	createAccountUsecase    app.CreateAccountUsecase
	getAccountUsecase       app.GetAccountUsecase
	listAccountsUsecase     app.ListAccountsUsecase
	createTransferUsecase   app.CreateTransferUsecase
	createUserUsecase       app.CreateUserUsecase
	loginUserUsecase        app.LoginUserUsecase
	renewAccessTokenUsecase app.RenewAccessTokenUsecase
}

func NewHandler(i *do.Injector) (*Handler, error) {
	r := do.MustInvoke[app.Repository](i)
	d := do.MustInvoke[app.Dispatcher](i)
	tm := do.MustInvoke[token.Manager](i)
	accessTokenDuration := do.MustInvokeNamed[time.Duration](i, "AccessTokenDuration")
	refreshTokenDuration := do.MustInvokeNamed[time.Duration](i, "RefreshTokenDuration")

	return &Handler{
		createAccountUsecase:    usecase.NewCreateAccountUsecase(r),
		getAccountUsecase:       usecase.NewGetAccountUsecase(r),
		listAccountsUsecase:     usecase.NewListAccountsUsecase(r),
		createTransferUsecase:   usecase.NewCreateTransferUsecase(r),
		createUserUsecase:       usecase.NewCreateUserUsecase(r, d),
		loginUserUsecase:        usecase.NewLoginUserUsecase(r, tm, accessTokenDuration, refreshTokenDuration),
		renewAccessTokenUsecase: usecase.NewRenewAccessTokenUsecase(r, tm, accessTokenDuration),
	}, nil
}
