package usecase

import (
	"time"

	"playground/internal/pkg/mail"
	"playground/internal/pkg/token"
	"playground/internal/wallet"
	"playground/internal/wallet/mq"
)

type Usecase struct {
	r                    wallet.Repository
	q                    mq.Producer
	tm                   token.Manager
	mailer               mail.Sender
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewUsecase(r wallet.Repository, q mq.Producer, mailer mail.Sender, tm token.Manager, accessTokenDuration, refreshTokenDuration time.Duration) wallet.Usecase {
	return &Usecase{
		r:                    r,
		q:                    q,
		mailer:               mailer,
		tm:                   tm,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}
