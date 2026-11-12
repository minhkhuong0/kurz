-- name: GetUrl :one
SELECT * FROM shorturls
WHERE code = $1;

-- name: ListUrls :many
SELECT * FROM shorturls
ORDER BY created_at;

-- name: CreateUrl :one
INSERT INTO shorturls (
  code, url, created_at, expires_at )
VALUES ($1, $2, $3, $3)
RETURNING *;

-- name: DeleteURL :exec
DELETE FROM shorturls
WHERE code = $1;
