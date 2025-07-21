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

-- name: UpdateUser :one
UPDATE users SET
    name = ?,
    email = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE id = ? RETURNING *;
