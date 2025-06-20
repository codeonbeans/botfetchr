-- name: GetUserAgent :one
SELECT *
FROM "account"."user_agents"
WHERE user_agent = $1;

-- name: CountUserAgents :one
SELECT COUNT(user_agent)
FROM "account"."user_agents"
WHERE (
  (good = sqlc.narg('good') OR sqlc.narg('good') IS NULL)
);

-- name: ListUserAgents :many
SELECT *
FROM "account"."user_agents"
WHERE (
  (good = sqlc.narg('good') OR sqlc.narg('good') IS NULL)
)
ORDER BY sqlc.arg('order_by') DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateUserAgent :one
INSERT INTO "account"."user_agents" (user_agent, good)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateUserAgent :one
UPDATE "account"."user_agents"
SET
  good = COALESCE(sqlc.narg('good'), good)
WHERE user_agent = $1
RETURNING *;

-- name: DeleteUserAgent :exec
DELETE FROM "account"."user_agents"
WHERE user_agent = $1;
