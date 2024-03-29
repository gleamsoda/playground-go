// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: users.sql

package gen

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :execlastid
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  ?, ?, ?, ?
)
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg *CreateUserParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

const getUser = `-- name: GetUser :one
SELECT username, hashed_password, full_name, email, password_changed_at, created_at, is_email_verified FROM users
WHERE username = ? LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (*User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsEmailVerified,
	)
	return &i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET
  hashed_password = ?,
  password_changed_at = ?,
  full_name = ?,
  email = ?,
  is_email_verified = ?
WHERE
  username = ?
`

type UpdateUserParams struct {
	HashedPassword    string    `json:"hashed_password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	IsEmailVerified   bool      `json:"is_email_verified"`
	Username          string    `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg *UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.FullName,
		arg.Email,
		arg.IsEmailVerified,
		arg.Username,
	)
	return err
}
