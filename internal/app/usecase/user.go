package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/morikuni/failure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"playground/internal/app"
	"playground/internal/pkg/apperr"
	"playground/internal/pkg/password"
	"playground/internal/pkg/token"
)

type (
	CreateUser struct {
		r app.Repository
		d app.Dispatcher
	}
	LoginUser struct {
		r                    app.Repository
		tm                   token.Manager
		accessTokenDuration  time.Duration
		refreshTokenDuration time.Duration
	}
	RenewAccessToken struct {
		r                   app.Repository
		tm                  token.Manager
		accessTokenDuration time.Duration
	}
	UpdateUser struct {
		r app.Repository
	}
)

func NewCreateUser(r app.Repository, d app.Dispatcher) *CreateUser {
	return &CreateUser{
		r: r,
		d: d,
	}
}

func (u *CreateUser) Execute(ctx context.Context, args *app.CreateUserParams) (*app.User, error) {
	hashedPassword, err := password.Hash(args.Password)
	if err != nil {
		return nil, err
	}

	var usr *app.User
	err = u.r.Transaction().Run(ctx, func(ctx context.Context, r app.Repository) error {
		var err error
		if usr, err = r.User().Create(ctx, app.NewUser(
			args.Username,
			hashedPassword,
			args.FullName,
			args.Email,
		)); err != nil {
			return err
		}
		if err := u.d.SendVerifyEmail(ctx, &app.SendVerifyEmailPayload{
			Username: usr.Username,
		}); err != nil {
			return err
		}
		return nil
	})
	return usr, err
}

func NewLoginUser(r app.Repository, tm token.Manager, accessTokenDuration time.Duration, refreshTokenDuration time.Duration) *LoginUser {
	return &LoginUser{
		r:                    r,
		tm:                   tm,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (u *LoginUser) Execute(ctx context.Context, args *app.LoginUserParams) (*app.LoginUserOutputParams, error) {
	usr, err := u.r.User().Get(ctx, args.Username)
	if err != nil {
		return nil, err
	}
	if err := password.Verify(args.Password, usr.HashedPassword); err != nil {
		return nil, failure.Translate(err, apperr.Unauthenticated)
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

	if err := u.r.User().CreateSession(ctx, app.NewSession(
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

func NewRenewAccessToken(r app.Repository, tm token.Manager, accessTokenDuration time.Duration) *RenewAccessToken {
	return &RenewAccessToken{
		r:                   r,
		tm:                  tm,
		accessTokenDuration: accessTokenDuration,
	}
}

func (u *RenewAccessToken) Execute(ctx context.Context, refreshToken string) (*app.RenewAccessTokenOutputParams, error) {
	rPayload, err := u.tm.Verify(refreshToken)
	if err != nil {
		return nil, err
	}

	sess, err := u.r.User().GetSession(ctx, rPayload.ID)
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

func NewUpdateUser(r app.Repository) *UpdateUser {
	return &UpdateUser{
		r: r,
	}
}

func (u *UpdateUser) Execute(ctx context.Context, args *app.UpdateUserParams) (*app.User, error) {
	if args.Username != args.ReqUsername {
		return nil, failure.Translate(fmt.Errorf("cannot update other user's info"), apperr.Unauthenticated)
	}

	usr, err := u.r.User().Get(ctx, args.Username)
	if err != nil {
		return nil, err
	}

	if args.FullName != nil {
		usr.FullName = *args.FullName
	}
	if args.Email != nil {
		usr.Email = *args.Email
	}
	if args.Password != nil {
		hashedPassword, err := password.Hash(*args.Password)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		usr.HashedPassword = hashedPassword
		usr.PasswordChangedAt = time.Now()
	}

	return u.r.User().Update(ctx, usr)
}
