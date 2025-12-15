-- name: AddObj :one
INSERT INTO mimix_obj (obj, obj_type, promote_date, obj_ver, lib, mimix_status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, obj, obj_type, promote_date, obj_ver, lib, mimix_status;

-- name: RemoveObj :exec
DELETE FROM mimix_obj
WHERE obj = $1;

-- name: UpdateObjStatus :exec
UPDATE mimix_obj
SET mimix_status = $2
WHERE obj = $1;