-- name: CreateCredit :one
INSERT INTO credit (
    client_id, bank_id, min_payment, max_payment, term_months, credit_type, status
) VALUES (
   $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetCredit :one
SELECT * FROM credit WHERE id = $1;

-- name: GetClientCredits :many
SELECT * FROM credit where client_id = $1;

-- name: UpdateCreditStatus :exec
UPDATE credit SET status = $1 WHERE id = $2;