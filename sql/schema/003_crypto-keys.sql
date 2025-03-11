-- +goose Up
CREATE TABLE crypto_keys (
    identity_key TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    signed_prekey TEXT NOT NULL,
    signed_key TEXT NOT NULL,
    onetime_prekey TEXT NOT NULL 
) ;

-- +goose Down
DROP TABLE crypto_keys ;
