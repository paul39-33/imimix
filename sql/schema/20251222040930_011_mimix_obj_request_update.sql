-- +goose Up
-- +goose StatementBegin
ALTER TABLE mimix_obj_req
ADD COLUMN promote_date DATE NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mimix_obj_req
DROP COLUMN promote_date;
-- +goose StatementEnd
