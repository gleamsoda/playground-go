-- name: CreateUser :execlastid
INSERT INTO users (
  username,
  full_name,
  email,
  hashed_password
) VALUES (
  ?, ?, ?, ?
);

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? LIMIT 1;