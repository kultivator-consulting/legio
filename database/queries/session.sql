-- name: GetSession :one
SELECT * FROM session
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetSessionsByAccountIdClientId :many
SELECT * FROM session
WHERE account_id = $1
  AND client_id = $2
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: GetSessionsByRefreshToken :many
SELECT * FROM session
WHERE refresh_token = $1
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: CountSessions :one
SELECT
    COUNT(id)
FROM session
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListSessions :many
SELECT * FROM session
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY
    CASE
        WHEN sqlc.arg(sort_by)::text = 'created' AND sqlc.arg(sort_order) = 'asc' THEN sqlc.arg(sort_by)::text END ASC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'created' AND sqlc.arg(sort_order) = 'desc' THEN sqlc.arg(sort_by)::text END DESC

OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateSession :one
INSERT INTO session (
    id,
    created,
    modified,
    deleted,
    account_id,
    client_id,
    client_agent,
    client_ip,
    client_bundle_id,
    access_token,
    refresh_token,
    access_token_expiry,
    refresh_token_expiry,
    is_device_app
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
         )
RETURNING *;

-- name: DeleteSession :exec
UPDATE session
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: DeleteSessionsByAccountIdClientId :exec
UPDATE session
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE account_id = $1
  AND client_id = $2;
