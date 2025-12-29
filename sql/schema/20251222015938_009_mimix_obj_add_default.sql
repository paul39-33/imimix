-- +goose Up
-- +goose StatementBegin
ALTER TABLE mimix_obj
ALTER COLUMN mimix_status SET DEFAULT 'unset';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mimix_obj
ALTER COLUMN mimix_status DROP DEFAULT;
-- +goose StatementEnd
