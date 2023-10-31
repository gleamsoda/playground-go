-- name: CreateVerifyEmail :execlastid
INSERT INTO verify_emails (
    username,
    email,
    secret_code
) VALUES (
    ?, ?, ?
);

-- name: GetVerifyEmail :one
SELECT * FROM verify_emails
WHERE id = ? LIMIT 1;

-- name: UpdateVerifyEmail :exec
UPDATE verify_emails
SET
    is_used = 1
WHERE
    id = ?
    AND secret_code = ?
    AND is_used = 0;
