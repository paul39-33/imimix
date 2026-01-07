-- name: CreateMimixObjReq :one
INSERT INTO mimix_obj_req (
    obj_name,
    requester,
    req_status,
    lib,
    obj_ver,
    obj_type,
    promote_date,
    developer
)
VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8
)
RETURNING id, obj_name, requester, req_status, lib, obj_ver, obj_type, promote_date, developer, created_at, updated_at;

-- name: UpdateMimixObjReqStatus :exec
UPDATE mimix_obj_req
SET req_status = $1
WHERE id = $2;

-- name: GetMimixObjReqByRequester :many
SELECT *
FROM mimix_obj_req
WHERE requester = $1;

-- name: GetMimixObjReq :many
SELECT *
FROM mimix_obj_req;

-- name: RemoveMimixObjReq :exec
DELETE FROM mimix_obj_req
WHERE id = $1
RETURNING obj_name;

-- name: GetMimixObjReqByID :one
SELECT *
FROM mimix_obj_req
WHERE id = $1;

-- name: CompleteMimixObjReq :exec
UPDATE mimix_obj_req
SET req_status = 'completed', updated_at = NOW()
WHERE id = $1;

-- name: UpdatePromoteStatus :exec
UPDATE mimix_obj_req
SET promote_status = $2
WHERE id = $1;

-- name: UpdateMimixObjReqInfo :one
UPDATE mimix_obj_req
SET obj_name = $2,
    lib = $3,
    obj_ver = $4,
    obj_type = $5,
    promote_date = $6,
    developer = $7,
    updated_at = NOW(),
    promote_status = $8,
    req_status = $9
WHERE id = $1
RETURNING id, obj_name, requester, req_status, lib, obj_ver, obj_type, promote_date, developer, created_at, updated_at, promote_status;

-- name: SearchMimixObjReq :many
SELECT * FROM mimix_obj_req
WHERE
    obj_name ILIKE '%' || $1 || '%'
 OR requester ILIKE '%' || $1 || '%'
 OR developer ILIKE '%' || $1 || '%'
 OR lib ILIKE '%' || $1 || '%'
ORDER BY updated_at DESC;

-- name: GetPendingObjReqByNameAndLib :one
SELECT *
FROM mimix_obj_req
WHERE obj_name = $1 AND lib = $2 AND req_status = 'pending';
