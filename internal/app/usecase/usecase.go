package usecase

import (
	"time"

	"github.com/samber/do"

	"playground/internal/app"
	"playground/internal/pkg/mail"
	"playground/internal/pkg/token"
)

type Usecase struct {
	r                    app.Repository
	d                    app.Dispatcher
	tm                   token.Manager
	mailer               mail.Sender
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewUsecase(i *do.Injector) (app.Usecase, error) {
	r := do.MustInvoke[app.Repository](i)
	d := do.MustInvoke[app.Dispatcher](i)
	tm := do.MustInvoke[token.Manager](i)
	mailer := do.MustInvoke[mail.Sender](i)
	accessTokenDuration := do.MustInvokeNamed[time.Duration](i, "AccessTokenDuration")
	refreshTokenDuration := do.MustInvokeNamed[time.Duration](i, "RefreshTokenDuration")

	return &Usecase{
		r:                    r,
		d:                    d,
		mailer:               mailer,
		tm:                   tm,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}, nil
}
