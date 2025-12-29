-- +goose Up
-- +goose StatementBegin
ALTER TABLE mimix_obj
ALTER COLUMN promote_date DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mimix_obj
ALTER COLUMN promote_date SET NOT NULL;
-- +goose StatementEnd
