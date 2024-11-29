-- name: GetPassProductTypeById :one
SELECT * FROM pass_product_type
WHERE id = $1
LIMIT 1;

-- name: CountPassProductTypes :one
SELECT
    COUNT(id)
FROM pass_product_type
WHERE pass_id = $1;

-- name: ListPassProductTypesByPassId :many
SELECT * FROM pass_product_type
WHERE pass_id = $1;

-- name: CreatePassProductType :one
INSERT INTO pass_product_type (
    id,
    created,
    modified,
    pass_id,
    product_type_id,
    category
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3
         )
RETURNING *;

-- name: DeletePassProductType :exec
DELETE FROM pass_product_type
WHERE id = $1;

-- name: DeletePassProductTypeByPassId :exec
DELETE FROM pass_product_type
WHERE pass_id = $1;

-- name: UpdatePassProductType :one
UPDATE pass_product_type
SET
    modified = now() AT TIME ZONE 'UTC',
    pass_id = $2,
    product_type_id = $3,
    category = $4
WHERE id = $1
RETURNING *;
