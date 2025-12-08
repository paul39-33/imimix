-- +goose Up
-- +goose StatementBegin
CREATE TABLE mimix_obj (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    obj TEXT NOT NULL,
    obj_type TEXT NOT NULL,
    promote_date DATE NOT NULL,
    lib TEXT NOT NULL,
    obj_ver TEXT NOT NULL,
    mimix_status TEXT NOT NULL,
    REFERENCES mimix_lib(lib)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE mimix;
-- +goose StatementEnd
