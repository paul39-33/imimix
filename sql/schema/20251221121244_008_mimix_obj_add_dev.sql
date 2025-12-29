-- +goose Up
-- +goose StatementBegin
ALTER TABLE mimix_obj
ADD developer TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mimix_obj
DROP COLUMN developer;
-- +goose StatementEnd
