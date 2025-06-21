-- name: GetSubscription :one
SELECT subscription.*
FROM "subscription"."subscriptions" subscription
WHERE id = $1;

-- name: CountSubscriptions :one
SELECT COUNT(id)
FROM "subscription"."subscriptions"
WHERE (
  (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL) AND
  (plan_id = sqlc.narg('plan_id') OR sqlc.narg('plan_id') IS NULL) AND
  (status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
  (start_date >= sqlc.narg('start_date_from') OR sqlc.narg('start_date_from') IS NULL) AND
  (start_date <= sqlc.narg('start_date_to') OR sqlc.narg('start_date_to') IS NULL) AND
  (end_date >= sqlc.narg('end_date_from') OR sqlc.narg('end_date_from') IS NULL) AND
  (end_date <= sqlc.narg('end_date_to') OR sqlc.narg('end_date_to') IS NULL) AND
  (cancel_at >= sqlc.narg('cancel_at_from') OR sqlc.narg('cancel_at_from') IS NULL) AND
  (cancel_at <= sqlc.narg('cancel_at_to') OR sqlc.narg('cancel_at_to') IS NULL)
);

-- name: ListSubscriptions :many
SELECT subscription.*
FROM "subscription"."subscriptions" subscription
WHERE (
  (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL) AND
  (plan_id = sqlc.narg('plan_id') OR sqlc.narg('plan_id') IS NULL) AND
  (status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
  (start_date >= sqlc.narg('start_date_from') OR sqlc.narg('start_date_from') IS NULL) AND
  (start_date <= sqlc.narg('start_date_to') OR sqlc.narg('start_date_to') IS NULL) AND
  (end_date >= sqlc.narg('end_date_from') OR sqlc.narg('end_date_from') IS NULL) AND
  (end_date <= sqlc.narg('end_date_to') OR sqlc.narg('end_date_to') IS NULL) AND
  (cancel_at >= sqlc.narg('cancel_at_from') OR sqlc.narg('cancel_at_from') IS NULL) AND
  (cancel_at <= sqlc.narg('cancel_at_to') OR sqlc.narg('cancel_at_to') IS NULL)
)
ORDER BY sqlc.arg('order_by') DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateSubscription :one
INSERT INTO "subscription"."subscriptions" (account_id, plan_id, status, start_date, end_date, cancel_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateSubscription :one
UPDATE "subscription"."subscriptions"
SET
  account_id = COALESCE(sqlc.narg('account_id'), account_id),
  plan_id = COALESCE(sqlc.narg('plan_id'), plan_id),
  status = COALESCE(sqlc.narg('status'), status),
  start_date = COALESCE(sqlc.narg('start_date'), start_date),
  end_date = CASE WHEN sqlc.arg('null_end_date')::boolean THEN NULL ELSE COALESCE(sqlc.narg('end_date'), end_date) END,
  cancel_at = CASE WHEN sqlc.arg('null_cancel_at')::boolean THEN NULL ELSE COALESCE(sqlc.narg('cancel_at'), cancel_at) END
WHERE id = $1
RETURNING *;

-- name: DeleteSubscription :exec
DELETE FROM "subscription"."subscriptions"
WHERE id = $1;
