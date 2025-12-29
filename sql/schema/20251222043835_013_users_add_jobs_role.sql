-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_job AS ENUM (
    'cmt',
    'dev',
    'dc',
    'user'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE user_job;
-- +goose StatementEnd
