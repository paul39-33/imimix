-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN job;

ALTER TABLE users
ADD job user_job NOT NULL DEFAULT 'user';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN job;
-- +goose StatementEnd
