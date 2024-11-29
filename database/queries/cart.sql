-- name: GetCartById :one
SELECT * FROM cart
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetCartByCartSessionId :one
SELECT * FROM cart
WHERE cart_session_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountCarts :one
SELECT
    COUNT(id)
FROM cart
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListAllCarts :many
SELECT * FROM cart
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY created;

-- name: CreateCart :one
INSERT INTO cart (
    id,
    created,
    modified,
    deleted,
    cart_session_id,
    domain_id,
    completed
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3
         )
RETURNING *;

-- name: DeleteCartById :exec
UPDATE cart
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateCartById :one
UPDATE cart
SET
    modified = now() AT TIME ZONE 'UTC',
    cart_session_id = $2,
    domain_id = $3,
    completed = $4
WHERE id = $1
RETURNING *;
