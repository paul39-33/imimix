-- name: CreateMimixObjReq :one
INSERT INTO mimix_obj_req (
    obj_name,
    requester,
    req_status,
    lib,
    obj_ver,
    obj_type,
    promote_date
)
VALUES (
    $1, $2, $3, $4,
    $5, $6, $7
)
RETURNING id, obj_name, requester, req_status, lib, obj_ver, obj_type, promote_date, created_at, updated_at;

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