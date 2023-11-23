package handler

import (
	"time"

	"github.com/samber/do"

	"playground/internal/app"
	"playground/internal/app/usecase"
	"playground/internal/pkg/token"
)

type Handler struct {
	createAccount    app.CreateAccountUsecase
	getAccount       app.GetAccountUsecase
	listAccounts     app.ListAccountsUsecase
	createTransfer   app.CreateTransferUsecase
	createUser       app.CreateUserUsecase
	loginUser        app.LoginUserUsecase
	renewAccessToken app.RenewAccessTokenUsecase
}

func NewHandler(i *do.Injector) (*Handler, error) {
	r := do.MustInvoke[app.Repository](i)
	d := do.MustInvoke[app.Dispatcher](i)
	tm := do.MustInvoke[token.Manager](i)
	accessTokenDuration := do.MustInvokeNamed[time.Duration](i, "AccessTokenDuration")
	refreshTokenDuration := do.MustInvokeNamed[time.Duration](i, "RefreshTokenDuration")

	return &Handler{
		createAccount:    usecase.NewCreateAccount(r),
		getAccount:       usecase.NewGetAccount(r),
		listAccounts:     usecase.NewListAccounts(r),
		createTransfer:   usecase.NewCreateTransfer(r),
		createUser:       usecase.NewCreateUser(r, d),
		loginUser:        usecase.NewLoginUser(r, tm, accessTokenDuration, refreshTokenDuration),
		renewAccessToken: usecase.NewRenewAccessToken(r, tm, accessTokenDuration),
	}, nil
}
