-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN IF NOT EXISTS role SMALLINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN IF EXISTS role;
-- +goose StatementEnd
