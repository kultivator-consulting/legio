-- name: GetPackagePricingById :one
SELECT * FROM package_pricing
WHERE id = $1
LIMIT 1;

-- name: GetPackagePricingByPackageId :many
SELECT * FROM package_pricing
WHERE package_id = $1;

-- name: GetValidPackagePricingByDate :many
SELECT * FROM package_pricing
WHERE package_id = $1
  AND start_date <= $2
  AND end_date >= $2;

-- name: CountPackagePricing :one
SELECT
    COUNT(id)
FROM package_pricing;

-- name: ListPackagePricingByPackageIdAsc :many
SELECT * FROM package_pricing
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackagePricingByPackageIdDesc :many
SELECT * FROM package_pricing
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackagePricing :one
INSERT INTO package_pricing (
    id,
    created,
    modified,
    package_id,
    supplement,
    adult_price,
    child_price,
    infant_price,
    start_date,
    end_date,
    specific_terms
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8
         )
RETURNING *;

-- name: UpdatePackagePricing :one
UPDATE package_pricing
SET
    modified = now() AT TIME ZONE 'UTC',
    package_id = $2,
    supplement = $3,
    adult_price = $4,
    child_price = $5,
    infant_price = $6,
    start_date = $7,
    end_date = $8,
    specific_terms = $9
WHERE id = $1
RETURNING *;

-- name: DeletePackagePricing :exec
DELETE FROM package_pricing
WHERE id = $1;

-- name: DeletePackagePricingByPackageId :exec
DELETE FROM package_pricing
WHERE package_id = $1;
