-- name: GetOperatorProductTypeById :one
SELECT * FROM operator_product_type
WHERE id = $1
LIMIT 1;

-- name: CountOperatorProductTypesByOperatorId :one
SELECT
    COUNT(id)
FROM operator_product_type
WHERE operator_id = $1;

-- name: ListOperatorProductTypesByOperatorId :many
SELECT * FROM operator_product_type
WHERE operator_id = $1;

-- name: CreateOperatorProductType :one
INSERT INTO operator_product_type (
    id,
    created,
    modified,
    operator_id,
    product_type_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: UpdateOperatorProductType :one
UPDATE operator_product_type
SET
    modified = now() AT TIME ZONE 'UTC',
    operator_id = $2,
    product_type_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteOperatorProductType :exec
DELETE FROM operator_product_type
WHERE id = $1;

-- name: DeleteOperatorProductTypeByOperatorId :exec
DELETE FROM operator_product_type
WHERE operator_id = $1;
