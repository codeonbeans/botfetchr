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
ORDER BY
  CASE WHEN sqlc.arg('order_by')::text = 'id_asc' THEN id END ASC,
  CASE WHEN sqlc.arg('order_by') = 'id_desc' THEN id END DESC,
  CASE WHEN sqlc.arg('order_by') = 'price_asc' THEN price END ASC,
  CASE WHEN sqlc.arg('order_by') = 'price_desc' THEN price END DESC,
  id ASC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreatePlan :one
INSERT INTO "subscription"."plans" (name, price, interval, description)
VALUES (
  sqlc.arg('name'),
  sqlc.arg('price'),
  sqlc.arg('interval'),
  sqlc.narg('description')
)
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
