-- +goose Up
-- +goose StatementBegin

-- 1. Create enum type if it doesn't exist
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

-- 2. Add column using the enum (if not exists)
ALTER TABLE mimix_obj
ADD COLUMN IF NOT EXISTS mimix_status mimix_status NOT NULL DEFAULT 'unset';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE mimix_obj
DROP COLUMN IF EXISTS mimix_status;

DROP TYPE IF EXISTS mimix_status;

-- +goose StatementEnd
