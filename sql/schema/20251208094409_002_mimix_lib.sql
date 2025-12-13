-- +goose Up
-- +goose StatementBegin
CREATE TABLE mimix_lib (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lib TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE mimix_lib;
-- +goose StatementEnd
