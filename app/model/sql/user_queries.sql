-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id;

-- name: CreateUser :one
INSERT INTO users ("first_name", "last_name", "email", "password", "last_login", "active", "role") VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: CreateUserWithId :one
INSERT INTO users ("id", "first_name", "last_name", "email", "password", "last_login", "active", "role") VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: DeleteUser :exec
UPDATE users SET
  deleted_at = now(),
  updated_at = now()
WHERE id = $1;

-- name: HardDeleteUser :exec
DELETE FROM users WHERE id = $1;

