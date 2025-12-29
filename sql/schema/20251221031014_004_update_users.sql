-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD job TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN job;
-- +goose StatementEnd
