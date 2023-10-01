package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/morikuni/failure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"playground/internal/pkg/apperr"
	"playground/internal/pkg/password"
	"playground/internal/wallet"
	"playground/internal/wallet/mq"
)

func (u *Usecase) CreateUser(ctx context.Context, args *wallet.CreateUserParams) (*wallet.User, error) {
	hashedPassword, err := password.Hash(args.Password)
	if err != nil {
		return nil, err
	}

	usr, err := u.r.CreateUser(ctx, wallet.NewUser(
		args.Username,
		hashedPassword,
		args.FullName,
		args.Email,
	))
	if err != nil {
		return nil, err
	}

	if err := u.q.SendVerifyEmail(ctx, &mq.SendVerifyEmailPayload{
		Username: usr.Username,
	}); err != nil {
		return nil, err
	}

	return usr, nil
}

func (u *Usecase) LoginUser(ctx context.Context, args *wallet.LoginUserParams) (*wallet.LoginUserOutputParams, error) {
	usr, err := u.r.GetUser(ctx, args.Username)
	if err != nil {
		return nil, err
	}
	if err := password.Verify(args.Password, usr.HashedPassword); err != nil {
		return nil, failure.Translate(err, apperr.Unauthorized)
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

	if err := u.r.CreateSession(ctx, wallet.NewSession(
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

	return &wallet.LoginUserOutputParams{
		SessionID:             rPayload.ID,
		AccessToken:           aToken,
		AccessTokenExpiresAt:  aPayload.ExpiredAt,
		RefreshToken:          rToken,
		RefreshTokenExpiresAt: rPayload.ExpiredAt,
		User:                  usr,
	}, nil
}

func (u *Usecase) RenewAccessToken(ctx context.Context, refreshToken string) (*wallet.RenewAccessTokenOutputParams, error) {
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

	resp := &wallet.RenewAccessTokenOutputParams{
		AccessToken:          aToken,
		AccessTokenExpiresAt: aPayload.ExpiredAt,
	}
	return resp, nil
}

func (u *Usecase) UpdateUser(ctx context.Context, args *wallet.UpdateUserParams) (*wallet.User, error) {
	if args.Username != args.ReqUsername {
		return nil, failure.Translate(fmt.Errorf("cannot update other user's info"), apperr.Unauthorized)
	}

	usr, err := u.r.GetUser(ctx, args.Username)
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

	return u.r.UpdateUser(ctx, usr)
}
