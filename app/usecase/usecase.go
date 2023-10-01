package usecase

import (
	"time"

	"playground/app"
	"playground/app/mq"
	"playground/pkg/mail"
	"playground/pkg/token"
)

type Usecase struct {
	r                    app.Repository
	q                    mq.Producer
	tm                   token.Manager
	mailer               mail.Sender
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewUsecase(r app.Repository, q mq.Producer, mailer mail.Sender, tm token.Manager, accessTokenDuration, refreshTokenDuration time.Duration) app.Usecase {
	return &Usecase{
		r:                    r,
		q:                    q,
		mailer:               mailer,
		tm:                   tm,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}
