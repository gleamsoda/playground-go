-- name: CreateEntry :execlastid
INSERT INTO entries (
  wallet_id,
  amount
) VALUES (
  ?, ?
);

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = ? LIMIT 1;
