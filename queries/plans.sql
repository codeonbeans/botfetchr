-- name: GetPlan :one
SELECT *
FROM "subscription"."plans"
WHERE id = $1;

-- name: CountPlans :one
SELECT COUNT(id)
FROM "subscription"."plans"
WHERE (
  (name = sqlc.narg('name') OR sqlc.narg('name') IS NULL) AND
  (price >= sqlc.narg('price_from') OR sqlc.narg('price_from') IS NULL) AND
  (price <= sqlc.narg('price_to') OR sqlc.narg('price_to') IS NULL) AND
  (interval = sqlc.narg('interval') OR sqlc.narg('interval') IS NULL) AND
  (description ILIKE '%' || sqlc.narg('description') || '%' OR sqlc.narg('description') IS NULL)
);

-- name: ListPlans :many
SELECT *
FROM "subscription"."plans"
WHERE (
  (name = sqlc.narg('name') OR sqlc.narg('name') IS NULL) AND
  (price >= sqlc.narg('price_from') OR sqlc.narg('price_from') IS NULL) AND
  (price <= sqlc.narg('price_to') OR sqlc.narg('price_to') IS NULL) AND
  (interval = sqlc.narg('interval') OR sqlc.narg('interval') IS NULL) AND
  (description ILIKE '%' || sqlc.narg('description') || '%' OR sqlc.narg('description') IS NULL)
)
ORDER BY sqlc.arg('order_by') DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreatePlan :one
INSERT INTO "subscription"."plans" (name, price, interval, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdatePlan :one
UPDATE "subscription"."plans"
SET
  name = COALESCE(sqlc.narg('name'), name),
  price = COALESCE(sqlc.narg('price'), price),
  interval = COALESCE(sqlc.narg('interval'), interval),
  description = CASE WHEN sqlc.arg('null_description')::boolean THEN NULL ELSE COALESCE(sqlc.narg('description'), description) END
WHERE id = $1
RETURNING *;

-- name: DeletePlan :exec
DELETE FROM "subscription"."plans"
WHERE id = $1;
