-- name: CreateUser :one
INSERT INTO users (
    id, 
    created_at, 
    updated_at, 
    email, 
    hashed_password, 
    identity_key, 
    signed_prekey, 
    signed_key
) VALUES (
    gen_random_uuid(), 
    NOW(), 
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING * ;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 ;

-- name: GetUserKeyPacket :one
SELECT identity_key, signed_prekey, signed_key FROM users 
WHERE id = $1 ;
