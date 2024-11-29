-- name: GetProductTypeDomainById :one
SELECT * FROM product_type_domain
WHERE id = $1
LIMIT 1;

-- name: CountProductTypeDomainsByProductTypeId :one
SELECT
    COUNT(id)
FROM product_type_domain
WHERE product_type_id = $1;

-- name: ListProductTypeDomainsByProductTypeId :many
SELECT * FROM product_type_domain
WHERE product_type_id = $1;

-- name: CreateProductTypeDomain :one
INSERT INTO product_type_domain (
    id,
    created,
    modified,
    product_type_id,
    domain_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: DeleteProductTypeDomain :exec
DELETE FROM product_type_domain
WHERE id = $1;

-- name: DeleteProductTypeDomainsByProductTypeId :exec
DELETE FROM product_type_domain
WHERE product_type_id = $1;

-- name: UpdateProductTypeDomain :one
UPDATE product_type_domain
SET
    modified = now() AT TIME ZONE 'UTC',
    product_type_id = $2,
    domain_id = $3
WHERE id = $1
RETURNING *;
