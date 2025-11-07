-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT refresh_tokens.*, users.* 
FROM refresh_tokens 
INNER JOIN users 
ON users.id = refresh_tokens.user_id
WHERE token = $1 AND revoked_at IS NULL;

-- name: SetRevokedAt :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = NOW()
WHERE token = $1;

