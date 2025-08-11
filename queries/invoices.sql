-- name: GetInvoice :one
SELECT *
FROM "subscription"."invoices"
WHERE id = $1;

-- name: CountInvoices :one
SELECT COUNT(id)
FROM "subscription"."invoices"
WHERE (
        (subscription_id = sqlc.narg('subscription_id') OR sqlc.narg('subscription_id') IS NULL) AND
        (amount >= sqlc.narg('amount_from') OR sqlc.narg('amount_from') IS NULL) AND
        (amount <= sqlc.narg('amount_to') OR sqlc.narg('amount_to') IS NULL) AND
        (issued_at >= sqlc.narg('issued_at_from') OR sqlc.narg('issued_at_from') IS NULL) AND
        (issued_at <= sqlc.narg('issued_at_to') OR sqlc.narg('issued_at_to') IS NULL) AND
        (paid = sqlc.narg('paid') OR sqlc.narg('paid') IS NULL)
        );

-- name: ListInvoices :many
SELECT *
FROM "subscription"."invoices"
WHERE (
        (subscription_id = sqlc.narg('subscription_id') OR sqlc.narg('subscription_id') IS NULL) AND
        (amount >= sqlc.narg('amount_from') OR sqlc.narg('amount_from') IS NULL) AND
        (amount <= sqlc.narg('amount_to') OR sqlc.narg('amount_to') IS NULL) AND
        (issued_at >= sqlc.narg('issued_at_from') OR sqlc.narg('issued_at_from') IS NULL) AND
        (issued_at <= sqlc.narg('issued_at_to') OR sqlc.narg('issued_at_to') IS NULL) AND
        (paid = sqlc.narg('paid') OR sqlc.narg('paid') IS NULL)
        )
ORDER BY CASE WHEN sqlc.arg('order_by')::text = 'id_asc' THEN id END ASC,
         CASE WHEN sqlc.arg('order_by') = 'id_desc' THEN id END DESC,
         CASE WHEN sqlc.arg('order_by') = 'issued_at_asc' THEN issued_at END ASC,
         CASE WHEN sqlc.arg('order_by') = 'issued_at_desc' THEN issued_at END DESC,
         issued_at DESC LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateInvoice :one
INSERT INTO "subscription"."invoices" (subscription_id, amount, issued_at, paid)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateInvoice :one
UPDATE "subscription"."invoices"
SET subscription_id = COALESCE(sqlc.narg('subscription_id'), subscription_id),
    amount          = COALESCE(sqlc.narg('amount'), amount),
    issued_at       = COALESCE(sqlc.narg('issued_at'), issued_at),
    paid            = COALESCE(sqlc.narg('paid'), paid)
WHERE id = $1 RETURNING *;

-- name: DeleteInvoice :exec
DELETE
FROM "subscription"."invoices"
WHERE id = $1;
