-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS forward_posts (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGSERIAL,
    from_chat_id BIGSERIAL,
    status SMALLINT NOT NULL DEFAULT 0,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_forward_posts_created_by ON posts (created_by);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS forward_posts;
-- +goose StatementEnd
