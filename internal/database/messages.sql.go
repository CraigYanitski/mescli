// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: messages.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createMessage = `-- name: CreateMessage :one
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
) RETURNING id, created_at, updated_at, user_id, sender_id, sender_identity_key, sender_ephemeral_key, message
`

type CreateMessageParams struct {
	UserID             uuid.UUID
	SenderID           uuid.UUID
	Message            string
	SenderIdentityKey  sql.NullString
	SenderEphemeralKey sql.NullString
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, createMessage,
		arg.UserID,
		arg.SenderID,
		arg.Message,
		arg.SenderIdentityKey,
		arg.SenderEphemeralKey,
	)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.SenderID,
		&i.SenderIdentityKey,
		&i.SenderEphemeralKey,
		&i.Message,
	)
	return i, err
}

const deleteMessage = `-- name: DeleteMessage :one
DELETE FROM messages 
WHERE id = $1 
RETURNING id, created_at, updated_at, user_id, sender_id, sender_identity_key, sender_ephemeral_key, message
`

func (q *Queries) DeleteMessage(ctx context.Context, id uuid.UUID) (Message, error) {
	row := q.db.QueryRowContext(ctx, deleteMessage, id)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.SenderID,
		&i.SenderIdentityKey,
		&i.SenderEphemeralKey,
		&i.Message,
	)
	return i, err
}

const getMessages = `-- name: GetMessages :many
SELECT id, created_at, updated_at, user_id, sender_id, sender_identity_key, sender_ephemeral_key, message FROM messages 
WHERE user_id = $1 
ORDER BY created_at
`

func (q *Queries) GetMessages(ctx context.Context, userID uuid.UUID) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessages, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.SenderID,
			&i.SenderIdentityKey,
			&i.SenderEphemeralKey,
			&i.Message,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
