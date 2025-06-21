-- name: GetResource :one
SELECT *
FROM "resource"."resources"
WHERE (
  id = sqlc.narg('id') OR
  resource = sqlc.narg('resource')
);

-- name: CountResources :one
SELECT COUNT(resource)
FROM "resource"."resources"
WHERE (
  (resource ILIKE '%' || sqlc.narg('resource') || '%' OR sqlc.narg('resource') IS NULL) AND
  (type = sqlc.narg('type') OR sqlc.narg('type') IS NULL) AND
  (attempts >= sqlc.narg('attempts_from') OR sqlc.narg('attempts_from') IS NULL) AND
  (attempts <= sqlc.narg('attempts_to') OR sqlc.narg('attempts_to') IS NULL) AND
  (failed = sqlc.narg('failed_from') OR sqlc.narg('failed_from') IS NULL) AND
  (failed <= sqlc.narg('failed_to') OR sqlc.narg('failed_to') IS NULL) AND
  (disabled = sqlc.narg('disabled') OR sqlc.narg('disabled') IS NULL)
);

-- name: ListResources :many
SELECT *
FROM "resource"."resources"
WHERE (
  (resource ILIKE '%' || sqlc.narg('resource') || '%' OR sqlc.narg('resource') IS NULL) AND
  (type = sqlc.narg('type') OR sqlc.narg('type') IS NULL) AND
  (attempts >= sqlc.narg('attempts_from') OR sqlc.narg('attempts_from') IS NULL) AND
  (attempts <= sqlc.narg('attempts_to') OR sqlc.narg('attempts_to') IS NULL) AND
  (failed = sqlc.narg('failed_from') OR sqlc.narg('failed_from') IS NULL) AND
  (failed <= sqlc.narg('failed_to') OR sqlc.narg('failed_to') IS NULL) AND
  (disabled = sqlc.narg('disabled') OR sqlc.narg('disabled') IS NULL)
)
ORDER BY sqlc.arg('order_by') DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateResource :one
INSERT INTO "resource"."resources" (resource, type, attempts, failed, disabled)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateResource :one
UPDATE "resource"."resources"
SET
  resource = COALESCE(sqlc.narg('resource'), resource),
  type = COALESCE(sqlc.narg('type'), type),
  attempts = COALESCE(sqlc.narg('attempts'), attempts),
  failed = COALESCE(sqlc.narg('failed'), failed),
  disabled = COALESCE(sqlc.narg('disabled'), disabled)
WHERE (
  id = sqlc.narg('id') OR
  resource = sqlc.narg('resource')
)
RETURNING *;

-- name: DeleteResource :exec
DELETE FROM "resource"."resources"
WHERE (
  id = sqlc.narg('id') OR
  resource = sqlc.narg('resource')
);
