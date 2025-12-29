-- +goose Up
-- +goose StatementBegin
CREATE TABLE mimix_obj_req (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    obj_name TEXT NOT NULL,
    requester TEXT NOT NULL REFERENCES users(username),
    req_status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE 'mimix_obj_req';
-- +goose StatementEnd
