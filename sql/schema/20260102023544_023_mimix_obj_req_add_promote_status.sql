-- +goose Up
-- +goose StatementBegin
ALTER TABLE mimix_obj_req ADD COLUMN promote_status promote_status DEFAULT 'in_progress';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mimix_obj_req DROP COLUMN promote_status;
-- +goose StatementEnd
