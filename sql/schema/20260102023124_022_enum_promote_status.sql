-- +goose Up
-- +goose StatementBegin
CREATE TYPE promote_status AS ENUM (
    'in_progress',
    'deployed'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE promote_status;
-- +goose StatementEnd
