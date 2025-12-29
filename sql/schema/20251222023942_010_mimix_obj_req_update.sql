-- +goose Up
-- +goose StatementBegin
ALTER TABLE mimix_obj_req
ADD COLUMN lib TEXT NOT NULL,
ADD COLUMN obj_ver TEXT NOT NULL,
ADD COLUMN obj_type TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mimix_obj_req
DROP COLUMN lib,
DROP COLUMN obj_ver,
DROP COLUMN obj_type;
-- +goose StatementEnd
