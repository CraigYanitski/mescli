// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CryptoKey struct {
	IdentityKey   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserID        uuid.UUID
	SignedPrekey  string
	SignedKey     string
	OnetimePrekey string
}

type Message struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	Message   string
}

type RefreshToken struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	Name           string
	HashedPassword string
	Initialised    bool
}
