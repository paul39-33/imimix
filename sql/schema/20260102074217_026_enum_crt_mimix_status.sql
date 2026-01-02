-- +goose Up
-- +goose StatementBegin

CREATE TYPE mimix_status AS ENUM (
  'unset',
  'done',
  'daftarkan',
  'tidak perlu daftar',
  'on progress'
);

ALTER TABLE mimix_obj
ALTER COLUMN mimix_status
TYPE mimix_status
USING mimix_status::mimix_status;

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
USING mimix_status::TEXT;

DROP TYPE mimix_status;

-- +goose StatementEnd
