-- name: GetPlanFeature :one
SELECT *
FROM "subscription"."plan_features"
WHERE id = $1;

-- name: CountPlanFeatures :one
SELECT COUNT(*)
FROM "subscription"."plan_features"
WHERE (
        (plan_id = sqlc.narg('plan_id') OR sqlc.narg('plan_id') IS NULL) AND
        (feature = sqlc.narg('feature') OR sqlc.narg('feature') IS NULL)
        );

-- name: ListPlanFeatures :many
SELECT *
FROM "subscription"."plan_features"
WHERE (
        (plan_id = sqlc.narg('plan_id') OR sqlc.narg('plan_id') IS NULL) AND
        (feature = sqlc.narg('feature') OR sqlc.narg('feature') IS NULL)
        )
ORDER BY id ASC LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreatePlanFeature :one
INSERT INTO "subscription"."plan_features" (plan_id, feature, "limit", days_to_reset)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdatePlanFeature :one
UPDATE "subscription"."plan_features"
SET plan_id       = COALESCE(sqlc.narg('plan_id'), plan_id),
    feature       = COALESCE(sqlc.narg('feature'), feature),
    "limit"       = COALESCE(sqlc.narg('limit'), "limit"),
    days_to_reset = COALESCE(sqlc.narg('days_to_reset'), days_to_reset)
WHERE id = $1 RETURNING *;

-- name: DeletePlanFeature :exec
DELETE
FROM "subscription"."plan_features"
WHERE id = $1;
