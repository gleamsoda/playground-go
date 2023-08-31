package usecase

import (
	"time"

	"playground/app"
	"playground/pkg/token"
)

type Usecase struct {
	r                    app.Repository
	tm                   token.Manager
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewUsecase(r app.Repository, tm token.Manager, accessTokenDuration, refreshTokenDuration time.Duration) app.Usecase {
	return &Usecase{
		r:                    r,
		tm:                   tm,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}
