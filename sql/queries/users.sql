-- name: CreateUser :one
INSERT INTO users (
    id, 
    created_at, 
    updated_at, 
    email, 
    name,
    hashed_password, 
    initialised
) VALUES (
    gen_random_uuid(), 
    NOW(), 
    NOW(),
    $1,
    $2,
    $3,
    false
) RETURNING * ;

-- name: GetUser :one
SELECT * FROM users 
WHERE id = $1 ;

-- name: GetUserByEmail :one
SELECT id FROM users 
WHERE email = $1 ;

-- name: UpdateUser :one
UPDATE users 
SET updated_at = NOW(),
    email = $2,
    name = $3,
    hashed_password = $4,
    initialised = $5
WHERE id = $1
RETURNING * ;

