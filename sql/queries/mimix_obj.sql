-- name: AddObj :one
INSERT INTO mimix_obj (obj, obj_type, promote_date, obj_ver, lib, lib_id, mimix_status, developer)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, obj, obj_type, promote_date, obj_ver, lib, lib_id, mimix_status, developer;

-- name: UpdateObjStatus :exec
UPDATE mimix_obj
SET mimix_status = $2
WHERE obj = $1
RETURNING obj, mimix_status;


-- name: GetObjByID :one
SELECT *
FROM mimix_obj
WHERE id = $1;


-- name: RemoveObjByID :exec
DELETE FROM mimix_obj
WHERE id = $1;

-- name: AddObjToObjReq :exec
INSERT INTO mimix_obj_req (
    obj_name,
    requester,
    req_status,
    created_at,
    updated_at,
    lib,
    obj_ver,
    obj_type,
    promote_date,
    developer,
    source_obj_id
)
SELECT
    o.obj,
    $2,
    $3,
    NOW(),
    NOW(),
    o.lib,
    o.obj_ver,
    o.obj_type,
    o.promote_date,
    o.developer,
    o.id          -- source obj id
FROM mimix_obj AS o
WHERE o.id = $1;

-- name: CompleteObjMimixStatus :exec
UPDATE mimix_obj
SET mimix_status = 'done'
WHERE id = $1;

-- name: UpdateObjInfo :one
UPDATE mimix_obj
SET
    obj           = $2,
    lib           = $3,
    obj_type      = $4,
    obj_ver       = $5,
    promote_date  = $6,
    mimix_status  = $7,
    developer     = $8,
    keterangan    = $9
WHERE id = $1
RETURNING *;

-- name: UpdateMimixStatus :exec
UPDATE mimix_obj
SET mimix_status = $2
WHERE id = $1;

-- name: GetMimixStatusByID :one
SELECT mimix_status
FROM mimix_obj
WHERE id = $1;

-- name: SearchMimixObj :many
SELECT *
FROM mimix_obj
WHERE
    obj ILIKE '%' || $1 || '%'
 OR lib ILIKE '%' || $1 || '%'
 OR developer ILIKE '%' || $1 || '%';