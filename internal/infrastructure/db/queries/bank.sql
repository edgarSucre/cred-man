-- name: CreateBank :one
INSERT INTO banks (name, type) VALUES ($1, $2) RETURNING *;


-- name: GetBank :one
SELECT * FROM banks
WHERE id = $1;