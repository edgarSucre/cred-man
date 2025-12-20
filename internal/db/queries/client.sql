-- name: GetClient :one
SELECT * FROM clients
WHERE id = $1;

-- name: CreateClient :one
INSERT INTO clients (
    full_name, email, birthdate, country
) VALUES (
    @full_name, @email, sqlc.narg('birthdate'), sqlc.narg('country')
) RETURNING *;
