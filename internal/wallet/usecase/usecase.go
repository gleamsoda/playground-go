package usecase

import (
	"time"

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

func NewUsecase(r wallet.Repository, q wallet.Dispatcher, mailer mail.Sender, tm token.Manager, accessTokenDuration, refreshTokenDuration time.Duration) wallet.Usecase {
	return &Usecase{
		r:                    r,
		d:                    q,
		mailer:               mailer,
		tm:                   tm,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}
