// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
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
) RETURNING id, created_at, updated_at, email, hashed_password, identity_key, signed_prekey, signed_key
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
	IdentityKey    string
	SignedPrekey   string
	SignedKey      string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.HashedPassword,
		arg.IdentityKey,
		arg.SignedPrekey,
		arg.SignedKey,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IdentityKey,
		&i.SignedPrekey,
		&i.SignedKey,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, identity_key, signed_prekey, signed_key FROM users 
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IdentityKey,
		&i.SignedPrekey,
		&i.SignedKey,
	)
	return i, err
}

const getUserKeyPacket = `-- name: GetUserKeyPacket :one
SELECT identity_key, signed_prekey, signed_key FROM users 
WHERE id = $1
`

type GetUserKeyPacketRow struct {
	IdentityKey  string
	SignedPrekey string
	SignedKey    string
}

func (q *Queries) GetUserKeyPacket(ctx context.Context, id uuid.UUID) (GetUserKeyPacketRow, error) {
	row := q.db.QueryRowContext(ctx, getUserKeyPacket, id)
	var i GetUserKeyPacketRow
	err := row.Scan(&i.IdentityKey, &i.SignedPrekey, &i.SignedKey)
	return i, err
}
