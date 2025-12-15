-- name: AddLib :one
INSERT INTO mimix_lib (lib)
VALUES ($1)
RETURNING id, lib;

-- name: DeleteLib :exec
DELETE FROM mimix_lib
WHERE lib = $1;