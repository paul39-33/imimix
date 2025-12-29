-- name: AddObj :one
INSERT INTO mimix_obj (obj, obj_type, promote_date, obj_ver, lib, lib_id, mimix_status, developer)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, obj, obj_type, promote_date, obj_ver, lib, lib_id, mimix_status, developer;

-- name: RemoveObj :exec
DELETE FROM mimix_obj
WHERE obj = $1;

-- name: UpdateObjStatus :exec
UPDATE mimix_obj
SET mimix_status = $2
WHERE obj = $1
RETURNING obj, mimix_status;

-- name: GetObjByName :one
SELECT *
FROM mimix_obj
WHERE obj = $1;

-- name: GetObjByDev :many
SELECT *
FROM mimix_obj
WHERE developer = $1;