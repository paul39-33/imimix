-- name: CreateUser :one
INSERT INTO users (username, hashed_password, job)
VALUES ($1, $2, $3)
RETURNING username, job, created_at, updated_at;

-- name: UserLogin :one
SELECT *
FROM users
WHERE username = $1;

-- name: GetUserByID :one
SELECT id, username, job, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT id, username, job, hashed_password, created_at, updated_at
FROM users
WHERE username = $1;