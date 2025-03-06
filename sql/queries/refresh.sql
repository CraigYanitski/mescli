-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING * ;

-- name: ResetRefreshTokenss :exec
DELETE FROM refresh_tokens ;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE user_id=$1
AND revoked_at IS NULL
AND expires_at > NOW() ;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
JOIN refresh_tokens ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL 
AND expires_at > NOW() ;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(),
    revoked_at = NOW()
WHERE token = $1 ;
