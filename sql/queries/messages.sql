-- name: CreateMessage :one
INSERT INTO messages (
    id,
    created_at,
    updated_at,
    user_id,
    sender_id,
    sender_identity_key,
    sender_ephemeral_key,
    message
) VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $4,
    $5,
    $3
) RETURNING * ;

-- name: GetMessages :many
SELECT * FROM messages 
WHERE user_id = $1 
ORDER BY created_at ;

-- name: DeleteMessage :one
DELETE FROM messages 
WHERE id = $1 
RETURNING * ;
