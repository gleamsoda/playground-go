-- name: CreateWallet :execlastid
INSERT INTO wallets (
  user_id,
  balance,
  currency
) VALUES (
  ?, ?, ?
);

-- name: GetWallet :one
SELECT * FROM wallets
WHERE id = ? LIMIT 1;

-- name: ListWallets :many
SELECT * FROM wallets
WHERE user_id = ?
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: DeleteWallet :exec
DELETE FROM wallets
WHERE id = ?;
