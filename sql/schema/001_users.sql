-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    hashed_password TEXT NOT NULL DEFAULT 'unset',
    initialised BOOLEAN NOT NULL DEFAULT false
) ;

-- +goose Down
DROP TABLE users ;
