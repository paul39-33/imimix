-- +goose Up
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN email TO username;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN username TO email;
-- +goose StatementEnd
