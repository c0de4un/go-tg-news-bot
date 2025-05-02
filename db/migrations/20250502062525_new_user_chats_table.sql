-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_chats (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    state INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_chats_user_id ON user_chats (user_id);
CREATE INDEX idx_user_chats_chat_id ON user_chats (chat_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_chats;
-- +goose StatementEnd
