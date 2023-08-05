-- name: CreateSession :exec
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
);

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = ? LIMIT 1;