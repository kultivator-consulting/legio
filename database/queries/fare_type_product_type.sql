-- name: GetFareTypeProductTypes :one
SELECT * FROM fare_type_product_type
WHERE id = $1
LIMIT 1;

-- name: CountFareTypeProductTypesByFareTypeId :one
SELECT
    COUNT(id)
FROM fare_type_product_type
WHERE fare_type_id = $1;

-- name: ListFareTypeProductTypesByFareTypeId :many
SELECT * FROM fare_type_product_type
WHERE fare_type_id = $1;

-- name: CreateFareTypeProductType :one
INSERT INTO fare_type_product_type (
    id,
    created,
    modified,
    fare_type_id,
    product_type_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: DeleteFareTypeProductType :exec
DELETE FROM fare_type_product_type
WHERE id = $1;

-- name: DeleteFareTypeProductTypesByFareTypeId :exec
DELETE FROM fare_type_product_type
WHERE fare_type_id = $1;

-- name: UpdateFareTypeProductType :one
UPDATE fare_type_product_type
SET
    modified = now() AT TIME ZONE 'UTC',
    fare_type_id = $2,
    product_type_id = $3
WHERE id = $1
RETURNING *;
