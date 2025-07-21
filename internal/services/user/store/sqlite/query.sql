-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
    id, name, email, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET
    name = ?,
    email = ?,
    updated_at = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
