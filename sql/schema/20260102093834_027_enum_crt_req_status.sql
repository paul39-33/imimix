-- +goose Up
-- +goose StatementBegin

-- 1. Drop column first (if it exists)
ALTER TABLE mimix_obj_req
DROP COLUMN IF EXISTS req_status;

-- 2. Drop enum type if it exists
DROP TYPE IF EXISTS req_status;

-- 3. Create enum type (simple, direct)
CREATE TYPE req_status AS ENUM (
    'pending',
    'completed'
);

-- 4. Re-add column using enum
ALTER TABLE mimix_obj_req
ADD COLUMN req_status req_status NOT NULL DEFAULT 'pending';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE mimix_obj_req
DROP COLUMN IF EXISTS req_status;

DROP TYPE IF EXISTS req_status;

-- +goose StatementEnd
