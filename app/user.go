package app

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewUser(username, fullName, email, hashedPassword string) *User {
	return &User{
		Username:       username,
		FullName:       fullName,
		Email:          email,
		HashedPassword: hashedPassword,
	}
}

type UserUsecase interface {
	Create(ctx context.Context, arg CreateUserInputParams) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Login(ctx context.Context, arg LoginUserInputParams) (*LoginUserOutputParams, error)
	RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenOutputParams, error)
}

type UserRepository interface {
	Create(ctx context.Context, u *User) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type CreateUserInputParams struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserInputParams struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserAgent string `json:"user_agent"`
	ClientIP  string `json:"client_ip"`
}

type LoginUserOutputParams struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  *User     `json:"user"`
}

type RenewAccessTokenOutputParams struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
