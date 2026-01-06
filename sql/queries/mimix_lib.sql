
-- name: GetMimixLibByName :one
SELECT id, lib
FROM mimix_lib
WHERE lib = $1;

-- name: CreateMimixLib :one
INSERT INTO mimix_lib (lib)
VALUES ($1)
RETURNING id, lib;

-- name: UpdateObjLibID :exec
UPDATE mimix_obj
SET lib_id = $2
WHERE id = $1;

