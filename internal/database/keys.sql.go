// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: keys.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createKeyPacket = `-- name: CreateKeyPacket :one
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
) RETURNING identity_key, created_at, updated_at, user_id, signed_prekey, signed_key, onetime_prekey
`

type CreateKeyPacketParams struct {
	UserID        uuid.UUID
	IdentityKey   string
	SignedPrekey  string
	SignedKey     string
	OnetimePrekey string
}

func (q *Queries) CreateKeyPacket(ctx context.Context, arg CreateKeyPacketParams) (CryptoKey, error) {
	row := q.db.QueryRowContext(ctx, createKeyPacket,
		arg.UserID,
		arg.IdentityKey,
		arg.SignedPrekey,
		arg.SignedKey,
		arg.OnetimePrekey,
	)
	var i CryptoKey
	err := row.Scan(
		&i.IdentityKey,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.SignedPrekey,
		&i.SignedKey,
		&i.OnetimePrekey,
	)
	return i, err
}

const getUserKeyPacket = `-- name: GetUserKeyPacket :one
SELECT identity_key, created_at, updated_at, user_id, signed_prekey, signed_key, onetime_prekey FROM crypto_keys 
WHERE user_id = $1
`

func (q *Queries) GetUserKeyPacket(ctx context.Context, userID uuid.UUID) (CryptoKey, error) {
	row := q.db.QueryRowContext(ctx, getUserKeyPacket, userID)
	var i CryptoKey
	err := row.Scan(
		&i.IdentityKey,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.SignedPrekey,
		&i.SignedKey,
		&i.OnetimePrekey,
	)
	return i, err
}

const updateKeyPacket = `-- name: UpdateKeyPacket :one
UPDATE crypto_keys 
SET updated_at = NOW(),
    identity_key = $2,
    signed_prekey = $3,
    signed_key = $4,
    onetime_prekey = $5
WHERE user_id = $1 
RETURNING identity_key, created_at, updated_at, user_id, signed_prekey, signed_key, onetime_prekey
`

type UpdateKeyPacketParams struct {
	UserID        uuid.UUID
	IdentityKey   string
	SignedPrekey  string
	SignedKey     string
	OnetimePrekey string
}

func (q *Queries) UpdateKeyPacket(ctx context.Context, arg UpdateKeyPacketParams) (CryptoKey, error) {
	row := q.db.QueryRowContext(ctx, updateKeyPacket,
		arg.UserID,
		arg.IdentityKey,
		arg.SignedPrekey,
		arg.SignedKey,
		arg.OnetimePrekey,
	)
	var i CryptoKey
	err := row.Scan(
		&i.IdentityKey,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.SignedPrekey,
		&i.SignedKey,
		&i.OnetimePrekey,
	)
	return i, err
}
