package usecase

import (
	"context"
	"fmt"
	"time"

	"playground/app"
	"playground/pkg/password"
	"playground/pkg/token"
)

type UserUsecase struct {
	ur                   app.UserRepository
	tm                   token.Manager
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewUserUsecase(ur app.UserRepository, tm token.Manager, accessTokenDuration, refreshTokenDuration time.Duration) app.UserUsecase {
	return &UserUsecase{
		ur:                   ur,
		tm:                   tm,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, args *app.CreateUserParams) (*app.User, error) {
	hashedPassword, err := password.HashPassword(args.Password)
	if err != nil {
		return nil, err
	}
	return u.ur.CreateUser(ctx, app.NewUser(
		args.Username,
		hashedPassword,
		args.FullName,
		args.Email,
	))
}

func (u *UserUsecase) Login(ctx context.Context, args *app.LoginUserParams) (*app.LoginUserOutputParams, error) {
	usr, err := u.ur.GetUser(ctx, args.Username)
	if err != nil {
		return nil, err
	}
	if err := password.CheckPassword(args.Password, usr.HashedPassword); err != nil {
		return nil, err
	}

	aToken, aPayload, err := u.tm.Create(
		usr.Username,
		u.accessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	rToken, rPayload, err := u.tm.Create(
		usr.Username,
		u.refreshTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	if err := u.ur.CreateSession(ctx, app.NewSession(
		rPayload.ID,
		usr.Username,
		rToken,
		args.UserAgent,
		args.ClientIP,
		false,
		rPayload.ExpiredAt,
	)); err != nil {
		return nil, err
	}

	return &app.LoginUserOutputParams{
		SessionID:             rPayload.ID,
		AccessToken:           aToken,
		AccessTokenExpiresAt:  aPayload.ExpiredAt,
		RefreshToken:          rToken,
		RefreshTokenExpiresAt: rPayload.ExpiredAt,
		User:                  usr,
	}, nil
}

func (u *UserUsecase) RenewAccessToken(ctx context.Context, refreshToken string) (*app.RenewAccessTokenOutputParams, error) {
	rPayload, err := u.tm.Verify(refreshToken)
	if err != nil {
		return nil, err
	}

	sess, err := u.ur.GetSession(ctx, rPayload.ID)
	if err != nil {
		return nil, err
	}
	if sess.IsBlocked {
		return nil, fmt.Errorf("blocked session")
	}
	if sess.Username != rPayload.Username {
		return nil, fmt.Errorf("incorrect session user")
	}
	if sess.RefreshToken != refreshToken {
		return nil, fmt.Errorf("mismatched session token")
	}
	if time.Now().After(sess.ExpiresAt) {
		return nil, fmt.Errorf("expired session")
	}

	aToken, aPayload, err := u.tm.Create(
		rPayload.Username,
		u.accessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	resp := &app.RenewAccessTokenOutputParams{
		AccessToken:          aToken,
		AccessTokenExpiresAt: aPayload.ExpiredAt,
	}
	return resp, nil
}
