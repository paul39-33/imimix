-- +goose Up
-- +goose StatementBegin
ALTER TABLE mimix_obj ADD COLUMN keterangan TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mimix_obj DROP COLUMN keterangan;
-- +goose StatementEnd
