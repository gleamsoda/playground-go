package domain

import (
	"context"
	"time"
)

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
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
	Create(ctx context.Context, arg CreateUserParams) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type CreateUserParams struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepository interface {
	Create(ctx context.Context, u *User) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
