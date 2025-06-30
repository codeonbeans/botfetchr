-- name: GetPlan :one
SELECT *
FROM "subscription"."plans"
WHERE id = $1;

-- name: CountPlans :one
SELECT COUNT(id)
FROM "subscription"."plans";

-- name: ListPlans :many
SELECT *
FROM "subscription"."plans"
ORDER BY id ASC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreatePlan :one
INSERT INTO "subscription"."plans" (id)
VALUES ($1)
RETURNING *;

-- name: UpdatePlan :one
UPDATE "subscription"."plans"
SET
  id = COALESCE(sqlc.narg('new_id'), id)
WHERE id = $1
RETURNING *;

-- name: DeletePlan :exec
DELETE FROM "subscription"."plans"
WHERE id = $1;
