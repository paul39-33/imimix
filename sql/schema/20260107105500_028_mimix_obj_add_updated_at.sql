-- +goose Up
ALTER TABLE mimix_obj
ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE mimix_obj
DROP COLUMN updated_at;
