-- +goose Up
-- +goose StatementBegin

-- 1. Remove default first
ALTER TABLE mimix_obj
ALTER COLUMN mimix_status DROP DEFAULT;

-- 2. Convert column to enum (NO CREATE TYPE)
ALTER TABLE mimix_obj
ALTER COLUMN mimix_status
TYPE mimix_status
USING mimix_status::mimix_status;

-- 3. Re-add default
ALTER TABLE mimix_obj
ALTER COLUMN mimix_status
SET DEFAULT 'unset';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE mimix_obj
ALTER COLUMN mimix_status DROP DEFAULT;

ALTER TABLE mimix_obj
ALTER COLUMN mimix_status
TYPE TEXT
USING mimix_status::text;

DROP TYPE mimix_status;

-- +goose StatementEnd
