package app

import (
	"context"
	"time"
)

type VerifyEmail struct {
	ID         int64
	Username   string
	Email      string
	SecretCode string
	IsUsed     bool
	ExpiredAt  time.Time
	CreatedAt  time.Time
}

func NewVerifyEmail(username, email, secretCode string) *VerifyEmail {
	return &VerifyEmail{
		Username:   username,
		Email:      email,
		SecretCode: secretCode,
	}
}

type (
	SendVerifyEmailUsecase interface {
		Execute(ctx context.Context, args *SendVerifyEmailPayload) (*VerifyEmail, error)
	}
	VerifyEmailUsecase interface {
		Execute(ctx context.Context, args *VerifyEmailParams) (*User, error)
	}
)

type VerifyEmailParams struct {
	EmailID    int64
	SecretCode string
}
