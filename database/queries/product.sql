-- name: GetProductById :one
SELECT * FROM product
WHERE id = $1
LIMIT 1;

-- name: ListProductByProductTypeId :many
SELECT * FROM product
WHERE product_type_id = $1
AND deleted > now() AT TIME ZONE 'UTC';

-- name: CountProducts :one
SELECT
    COUNT(id)
FROM product
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListProductsAsc :many
SELECT * FROM product
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListProductsDesc :many
SELECT * FROM product
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateProduct :one
INSERT INTO product (
    id,
    created,
    modified,
    deleted,
    product_type_id,
    operator_id,
    start_location_id,
    end_location_id,
    name,
    start_place,
    end_place,
    operator_code,
    instructions,
    notes,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
         )
RETURNING *;

-- name: DeleteProduct :exec
UPDATE service
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateProduct :one
UPDATE product
SET
    modified = now() AT TIME ZONE 'UTC',
    product_type_id = $2,
    operator_id = $3,
    start_location_id = $4,
    end_location_id = $5,
    name = $6,
    start_place = $7,
    end_place = $8,
    operator_code = $9,
    instructions = $10,
    notes = $11,
    is_active = $12
WHERE id = $1
RETURNING *;
