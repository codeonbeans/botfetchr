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
ORDER BY
  CASE WHEN sqlc.arg('order_by')::text = 'id_asc' THEN id END ASC,
  CASE WHEN sqlc.arg('order_by') = 'id_desc' THEN id END DESC,
  CASE WHEN sqlc.arg('order_by') = 'start_date_asc' THEN start_date END ASC,
  CASE WHEN sqlc.arg('order_by') = 'start_date_desc' THEN start_date END DESC,
  CASE WHEN sqlc.arg('order_by') = 'end_date_asc' THEN end_date END ASC,
  CASE WHEN sqlc.arg('order_by') = 'end_date_desc' THEN end_date END DESC,
  CASE WHEN sqlc.arg('order_by') = 'cancel_at_asc' THEN cancel_at END ASC,
  CASE WHEN sqlc.arg('order_by') = 'cancel_at_desc' THEN cancel_at END DESC,
  start_date DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateSubscription :one
INSERT INTO "subscription"."subscriptions" (account_id, plan_id, status, start_date, end_date, cancel_at)
VALUES (
  sqlc.arg('account_id'),
  sqlc.arg('plan_id'),
  sqlc.arg('status'),
  sqlc.arg('start_date'),
  sqlc.narg('end_date'),
  sqlc.narg('cancel_at')
)
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
