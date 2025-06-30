-- name: GetPlanPrice :one
SELECT *
FROM "subscription"."plan_prices"
WHERE id = $1;

-- name: CountPlanPrices :one
SELECT COUNT(*)
FROM "subscription"."plan_prices"
WHERE (
  (plan_id = sqlc.narg('plan_id') OR sqlc.narg('plan_id') IS NULL) AND
  (price >= sqlc.narg('price_from') OR sqlc.narg('price_from') IS NULL) AND
  (price <= sqlc.narg('price_to') OR sqlc.narg('price_to') IS NULL) AND
  (interval = sqlc.narg('interval') OR sqlc.narg('interval') IS NULL)
);

-- name: ListPlanPrices :many
SELECT *
FROM "subscription"."plan_prices"
WHERE (
  (plan_id = sqlc.narg('plan_id') OR sqlc.narg('plan_id') IS NULL) AND
  (price >= sqlc.narg('price_from') OR sqlc.narg('price_from') IS NULL) AND
  (price <= sqlc.narg('price_to') OR sqlc.narg('price_to') IS NULL) AND
  (interval = sqlc.narg('interval') OR sqlc.narg('interval') IS NULL)
)
ORDER BY id ASC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreatePlanPrice :one
INSERT INTO "subscription"."plan_prices" (plan_id, price, interval)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdatePlanPrice :one
UPDATE "subscription"."plan_prices"
SET
  plan_id = COALESCE(sqlc.narg('plan_id'), plan_id),
  price = COALESCE(sqlc.narg('price'), price),
  interval = COALESCE(sqlc.narg('interval'), interval)
WHERE id = $1
RETURNING *;

-- name: DeletePlanPrice :exec
DELETE FROM "subscription"."plan_prices"
WHERE id = $1;
