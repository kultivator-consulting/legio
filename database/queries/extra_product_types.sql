-- name: GetExtraProductTypeById :one
SELECT * FROM extra_product_type
WHERE id = $1
LIMIT 1;

-- name: CountExtraProductTypesByExtraId :one
SELECT
    COUNT(id)
FROM extra_product_type
WHERE extra_id = $1;

-- name: CountExtraProductTypesByProductTypeId :one
SELECT
    COUNT(id)
FROM extra_product_type
WHERE product_type_id = $1;

-- name: ListExtraProductTypesByExtraId :many
SELECT * FROM extra_product_type
WHERE extra_id = $1;

-- name: ListExtraProductTypesByProductTypeId :many
SELECT * FROM extra_product_type
WHERE product_type_id = $1;

-- name: CreateExtraProductType :one
INSERT INTO extra_product_type (
    id,
    created,
    modified,
    extra_id,
    product_type_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: DeleteExtraProductType :exec
DELETE FROM extra_product_type
WHERE id = $1;

-- name: DeleteExtraProductTypeByExtraId :exec
DELETE FROM extra_product_type
WHERE extra_id = $1;

-- name: UpdateExtraProductType :one
UPDATE extra_product_type
SET
    modified = now() AT TIME ZONE 'UTC',
    extra_id = $2,
    product_type_id = $3
WHERE id = $1
RETURNING *;
