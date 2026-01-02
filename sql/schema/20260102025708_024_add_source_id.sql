-- +goose Up
-- +goose StatementBegin
ALTER TABLE mimix_obj_req ADD COLUMN source_obj_id UUID NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mimix_obj_req DROP COLUMN source_obj_id;
-- +goose StatementEnd
