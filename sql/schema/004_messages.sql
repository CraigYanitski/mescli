-- +goose Up
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users,
    message TEXT NOT NULL
) ;

-- +goose Down
DROP TABLE messages ;
