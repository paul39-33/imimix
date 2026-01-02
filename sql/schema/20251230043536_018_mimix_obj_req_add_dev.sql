-- +goose Up
-- +goose StatementBegin
ALTER TABLE mimix_obj_req
ADD COLUMN developer TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mimix_obj_req DROP COLUMN developer;
-- +goose StatementEnd
