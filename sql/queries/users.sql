-- name: CreateUser :one
INSERT INTO users (email, pass)
VALUES ($1, $2)
RETURNING email, created_at, updated_at;

-- name: UserLogin :one
SELECT *
FROM users
WHERE email = $1;