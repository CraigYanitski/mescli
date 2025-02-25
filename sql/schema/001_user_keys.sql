-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL DEFAULT 'unset',
    identity_key TEXT NOT NULL,
    signed_prekey TEXT NOT NULL,
    signed_key TEXT NOT NULL
) ;

-- +goose Down
DROP TABLE users ;
