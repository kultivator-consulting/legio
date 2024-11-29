-- name: GetOperatorDomainById :one
SELECT * FROM operator_domain
WHERE id = $1
LIMIT 1;

-- name: GetProductByName :one
SELECT * FROM product
WHERE name = $1
LIMIT 1;

-- name: CountOperatorDomainsByOperatorId :one
SELECT
    COUNT(id)
FROM operator_domain
WHERE operator_id = $1;

-- name: ListOperatorDomainsByOperatorId :many
SELECT * FROM operator_domain
WHERE operator_id = $1;

-- name: CreateOperatorDomain :one
INSERT INTO operator_domain (
    id,
    created,
    modified,
    operator_id,
    domain_id
) VALUES (
    uuid_generate_v4(),
    now() AT TIME ZONE 'UTC',
    now() AT TIME ZONE 'UTC',
    $1, $2
)
RETURNING *;

-- name: UpdateOperatorDomain :one
UPDATE operator_domain
SET
    modified = now() AT TIME ZONE 'UTC',
    operator_id = $2,
    domain_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteOperatorDomain :exec
DELETE FROM operator_domain
WHERE id = $1;

-- name: DeleteOperatorDomainsByOperatorId :exec
DELETE FROM operator_domain
WHERE operator_id = $1;
