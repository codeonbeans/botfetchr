-- name: GetAccountUsage :one
SELECT *
FROM "account"."usage"
WHERE id = $1;

-- name: CountAccountUsages :one
SELECT COUNT(*)
FROM "account"."usage"
WHERE (
        (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL) AND
        (feature = sqlc.narg('feature') OR sqlc.narg('feature') IS NULL)
        );

-- name: ListAccountUsages :many
SELECT *
FROM "account"."usage"
WHERE (
        (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL) AND
        (feature = sqlc.narg('feature') OR sqlc.narg('feature') IS NULL)
        )
ORDER BY id ASC LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateAccountUsage :one
INSERT INTO "account"."usage" (account_id, feature, usage)
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateAccountUsage :one
UPDATE "account"."usage"
SET account_id = COALESCE(sqlc.narg('account_id'), account_id),
    feature    = COALESCE(sqlc.narg('feature'), feature),
    usage      = COALESCE(sqlc.narg('usage'), usage),
    reset_at   = COALESCE(sqlc.narg('reset_at'), reset_at)
WHERE id = $1 RETURNING *;

-- name: DeleteAccountUsage :exec
DELETE
FROM "account"."usage"
WHERE id = $1;
