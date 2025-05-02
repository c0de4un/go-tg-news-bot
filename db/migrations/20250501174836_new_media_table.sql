-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS media (
    id BIGSERIAL PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    local_src VARCHAR(255) NOT NULL,
    extension VARCHAR(8) NOT NULL,
    mime_type VARCHAR(32) NOT NULL,
    length INT NOT NULL,
    uploaded_by BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_post_id ON media (post_id);
CREATE INDEX idx_uploaded_by ON media (uploaded_by);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS media;
-- +goose StatementEnd
