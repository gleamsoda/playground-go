-- name: CreateUser :execlastid
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  ?, ?, ?, ?
);

-- name: GetUser :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET
  hashed_password = sqlc.arg(hashed_password),
  password_changed_at = sqlc.arg(password_changed_at),
  full_name = sqlc.arg(full_name),
  email = sqlc.arg(email),
  is_email_verified = sqlc.arg(is_email_verified)
WHERE
  username = sqlc.arg(username);
