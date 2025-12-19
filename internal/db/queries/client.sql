-- name: GetClient :one
SELECT * FROM clients
WHERE id = $1;

-- name: CreateClient :one
INSERT INTO clients (
    full_name, email, birthdate, country
) VALUES (
    $1, $2, $3, $4
) RETURNING *;
