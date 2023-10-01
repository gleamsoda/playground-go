package app

import "time"

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

type VerifyEmailParams struct {
	EmailID    int64
	SecretCode string
}
