-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
RENAME COLUMN pass TO hashed_password;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
RENAME COLUMN hashed_password TO pass;
-- +goose StatementEnd
