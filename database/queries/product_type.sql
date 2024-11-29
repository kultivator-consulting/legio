-- name: GetProductType :one
SELECT * FROM product_type
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetProductTypeById :one
SELECT * FROM product_type
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetProductTypeByName :one
SELECT * FROM product_type
WHERE name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: GetProductTypeBySlug :one
SELECT * FROM product_type
WHERE slug = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: CountProductTypes :one
SELECT
    COUNT(id)
FROM product_type
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListProductTypesAsc :many
SELECT * FROM product_type
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListProductTypesDesc :many
SELECT * FROM product_type
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListProductTypesByPassType :many
SELECT * FROM product_type
WHERE is_pass_type = TRUE
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: ListProductTypesByPassService :many
SELECT * FROM product_type
WHERE is_pass_service = TRUE
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: CreateProductType :one
INSERT INTO product_type (
    id,
    created,
    modified,
    deleted,
    name,
    priority,
    slug,
    terms,
    is_pass_type,
    is_pass_service,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7
         )
RETURNING *;

-- name: CreateProductTypeAndReturnId :one
INSERT INTO product_type (
    id,
    created,
    modified,
    deleted,
    name,
    priority,
    slug,
    terms,
    is_pass_type,
    is_pass_service,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7
         )
RETURNING id;

-- name: DeleteProductType :exec
UPDATE product_type
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateProductType :one
UPDATE product_type
SET
    modified = now() AT TIME ZONE 'UTC',
    name = $2,
    priority = $3,
    slug = $4,
    terms = $5,
    is_pass_type = $6,
    is_pass_service = $7,
    is_active = $8
WHERE id = $1
RETURNING *;
