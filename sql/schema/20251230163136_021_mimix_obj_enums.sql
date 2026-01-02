-- +goose Up
-- +goose StatementBegin

-- 1. Create enum if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'mimix_status'
    ) THEN
        CREATE TYPE mimix_status AS ENUM (
            'unset',
            'done',
            'daftarkan',
            'tidak perlu daftar',
            'on progress'
        );
    END IF;
END$$;

-- 2. DROP default first (CRITICAL)
ALTER TABLE mimix_obj
ALTER COLUMN mimix_status DROP DEFAULT;

-- 3. Convert TEXT → ENUM
ALTER TABLE mimix_obj
ALTER COLUMN mimix_status
TYPE mimix_status
USING mimix_status::mimix_status;

-- 4. Re-add default (now enum-aware)
ALTER TABLE mimix_obj
ALTER COLUMN mimix_status
SET DEFAULT 'unset'::mimix_status;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE mimix_obj
ALTER COLUMN mimix_status DROP DEFAULT;

ALTER TABLE mimix_obj
ALTER COLUMN mimix_status
TYPE TEXT
USING mimix_status::TEXT;

DROP TYPE IF EXISTS mimix_status;

-- +goose StatementEnd
