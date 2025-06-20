-- name: GetAccountTelegram :one
SELECT *
FROM "account"."telegrams"
WHERE id = $1;

-- name: CountAccountTelegrams :one
SELECT COUNT(id)
FROM "account"."telegrams"
WHERE (
  (telegram_id = sqlc.narg('telegram_id') OR sqlc.narg('telegram_id') IS NULL) AND
  (is_bot = sqlc.narg('is_bot') OR sqlc.narg('is_bot') IS NULL) AND
  (first_name = sqlc.narg('first_name') OR sqlc.narg('first_name') IS NULL) AND
  (last_name = sqlc.narg('last_name') OR sqlc.narg('last_name') IS NULL) AND
  (username = sqlc.narg('username') OR sqlc.narg('username') IS NULL) AND
  (language_code = sqlc.narg('language_code') OR sqlc.narg('language_code') IS NULL) AND
  (photo_url = sqlc.narg('photo_url') OR sqlc.narg('photo_url') IS NULL) AND
  (is_premium = sqlc.narg('is_premium') OR sqlc.narg('is_premium') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
);

-- name: ListAccountTelegrams :many
SELECT *
FROM "account"."telegrams"
WHERE (
  (telegram_id = sqlc.narg('telegram_id') OR sqlc.narg('telegram_id') IS NULL) AND
  (is_bot = sqlc.narg('is_bot') OR sqlc.narg('is_bot') IS NULL) AND
  (first_name = sqlc.narg('first_name') OR sqlc.narg('first_name') IS NULL) AND
  (last_name = sqlc.narg('last_name') OR sqlc.narg('last_name') IS NULL) AND
  (username = sqlc.narg('username') OR sqlc.narg('username') IS NULL) AND
  (language_code = sqlc.narg('language_code') OR sqlc.narg('language_code') IS NULL) AND
  (photo_url = sqlc.narg('photo_url') OR sqlc.narg('photo_url') IS NULL) AND
  (is_premium = sqlc.narg('is_premium') OR sqlc.narg('is_premium') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
)
ORDER BY sqlc.arg('order_by') DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateAccountTelegram :one
INSERT INTO "account"."telegrams" (telegram_id, is_bot, first_name, last_name, username, language_code, photo_url, is_premium)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateAccountTelegram :one
UPDATE "account"."telegrams"
SET
  telegram_id = COALESCE(sqlc.narg('telegram_id'), telegram_id),
  is_bot = COALESCE(sqlc.narg('is_bot'), is_bot),
  first_name = COALESCE(sqlc.narg('first_name'), first_name),
  last_name = COALESCE(sqlc.narg('last_name'), last_name),
  username = CASE WHEN sqlc.arg('null_username')::boolean THEN NULL ELSE COALESCE(sqlc.narg('username'), username) END,
  language_code = COALESCE(sqlc.narg('language_code'), language_code),
  photo_url = CASE WHEN sqlc.arg('null_photo_url')::boolean THEN NULL ELSE COALESCE(sqlc.narg('photo_url'), photo_url) END,
  is_premium = COALESCE(sqlc.narg('is_premium'), is_premium),
  created_at = COALESCE(sqlc.narg('created_at'), created_at)
WHERE id = $1
RETURNING *;

-- name: DeleteAccountTelegram :exec
DELETE FROM "account"."telegrams"
WHERE id = $1;
