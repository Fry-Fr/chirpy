-- +goose Up
CREATE TABLE refresh_tokens (
    token VARCHAR(255) NOT NULL PRIMARY KEY,
    user_id UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP DEFAULT NULL
);

-- +goose Down
DROP TABLE refresh_tokens;
