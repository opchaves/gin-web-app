-- name: GetWorkspaceByID :one
SELECT * FROM workspaces WHERE id = $1;

-- name: GetUserWorkspaces :many
SELECT * FROM workspaces WHERE user_id = $1;

-- name: CreateWorkspace :one
INSERT INTO workspaces ("name", "description", "currency", "language", "user_id") VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateWorkspace :exec
UPDATE workspaces SET "name" = $2, "description" = $3, "currency" = $4, "language" = $5 WHERE id = $1;

-- name: DeleteWorkspaces :exec
DELETE FROM workspaces;