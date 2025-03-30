-- name: CreateKeyPacket :one
INSERT INTO crypto_keys (
    identity_key,
    created_at,
    updated_at,
    user_id,
    signed_prekey,
    signed_key,
    onetime_prekey
) VALUES(
    $2,
    NOW(),
    NOW(),
    $1,
    $3,
    $4,
    $5
) RETURNING * ;

-- name: GetUserKeyPacket :one
SELECT * FROM crypto_keys 
WHERE user_id = $1 ;

-- name: GetUserIdentityKey :one
SELECT identity_key FROM crypto_keys 
WHERE user_id = $1 ;

-- name: UpdateKeyPacket :one
UPDATE crypto_keys 
SET updated_at = NOW(),
    identity_key = $2,
    signed_prekey = $3,
    signed_key = $4,
    onetime_prekey = $5
WHERE user_id = $1 
RETURNING * ;
