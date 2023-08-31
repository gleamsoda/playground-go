package usecase

import (
	"context"
	"fmt"
	"time"

	"playground/app"
	"playground/pkg/password"
)

func (u *Usecase) CreateUser(ctx context.Context, args *app.CreateUserParams) (*app.User, error) {
	hashedPassword, err := password.HashPassword(args.Password)
	if err != nil {
		return nil, err
	}
	return u.r.CreateUser(ctx, app.NewUser(
		args.Username,
		hashedPassword,
		args.FullName,
		args.Email,
	))
}

func (u *Usecase) LoginUser(ctx context.Context, args *app.LoginUserParams) (*app.LoginUserOutputParams, error) {
	usr, err := u.r.GetUser(ctx, args.Username)
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

	if err := u.r.CreateSession(ctx, app.NewSession(
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

func (u *Usecase) RenewAccessToken(ctx context.Context, refreshToken string) (*app.RenewAccessTokenOutputParams, error) {
	rPayload, err := u.tm.Verify(refreshToken)
	if err != nil {
		return nil, err
	}

	sess, err := u.r.GetSession(ctx, rPayload.ID)
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
