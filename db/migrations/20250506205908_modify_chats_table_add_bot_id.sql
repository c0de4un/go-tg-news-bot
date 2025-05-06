-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_chats
ADD COLUMN IF NOT EXISTS bot_type SMALLINT NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS bot_id BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN IF EXISTS role;
-- +goose StatementEnd
