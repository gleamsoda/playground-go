package usecase

import (
	"time"

	"github.com/samber/do"

	"playground/internal/pkg/mail"
	"playground/internal/pkg/token"
	"playground/internal/wallet"
)

type Usecase struct {
	r                    wallet.Repository
	d                    wallet.Dispatcher
	tm                   token.Manager
	mailer               mail.Sender
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewUsecase(i *do.Injector) (wallet.Usecase, error) {
	r := do.MustInvoke[wallet.Repository](i)
	d := do.MustInvoke[wallet.Dispatcher](i)
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
