-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users ORDER BY id;

-- name: InsertUser :one
INSERT INTO users (
  "first_name", "last_name", "email", "password", "last_login", "active"
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;
