//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package app

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"-"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	IsEmailVerified   bool      `json:"is_email_verified"`
}

func NewUser(username, hashedPassword, fullName, email string) *User {
	return &User{
		Username:       username,
		FullName:       fullName,
		Email:          email,
		HashedPassword: hashedPassword,
	}
}

type (
	CreateUserUsecase interface {
		Execute(ctx context.Context, args *CreateUserParams) (*User, error)
	}
	CreateUserParams struct {
		Username string `json:"username" binding:"required,alphanum"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}
	LoginUserUsecase interface {
		Execute(ctx context.Context, args *LoginUserParams) (*LoginUserOutputParams, error)
	}
	LoginUserParams struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		UserAgent string `json:"user_agent"`
		ClientIP  string `json:"client_ip"`
	}
	LoginUserOutputParams struct {
		SessionID             uuid.UUID `json:"session_id"`
		AccessToken           string    `json:"access_token"`
		AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
		RefreshToken          string    `json:"refresh_token"`
		RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
		User                  *User     `json:"user"`
	}
	RenewAccessTokenUsecase interface {
		Execute(ctx context.Context, refreshToken string) (*RenewAccessTokenOutputParams, error)
	}
	RenewAccessTokenOutputParams struct {
		AccessToken          string    `json:"access_token"`
		AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	}
	UpdateUserUsecase interface {
		Execute(ctx context.Context, args *UpdateUserParams) (*User, error)
	}
	UpdateUserParams struct {
		ReqUsername string  `binding:"required,alphanum"`
		Username    string  `json:"username" binding:"required,alphanum"`
		Password    *string `json:"password" binding:"min=6"`
		FullName    *string `json:"full_name" binding:""`
		Email       *string `json:"email" binding:"email"`
	}
	UserRepository interface {
		Create(ctx context.Context, args *User) (*User, error)
		Get(ctx context.Context, username string) (*User, error)
		Update(ctx context.Context, args *User) (*User, error)
		UpdateEmailVerified(ctx context.Context, args *VerifyEmailParams) (*User, *VerifyEmail, error)
		CreateSession(ctx context.Context, args *Session) error
		GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
		CreateVerifyEmail(ctx context.Context, args *VerifyEmail) (*VerifyEmail, error)
	}
)
