-- name: GetPackageTermsConditionsById :one
SELECT * FROM package_terms_conditions
WHERE id = $1
LIMIT 1;

-- name: GetPackageTermsConditionsByPackageId :many
SELECT * FROM package_terms_conditions
WHERE package_id = $1;

-- name: CountPackageTermsConditionsByPackageId :one
SELECT
    COUNT(id)
FROM package_terms_conditions
WHERE package_id = $1;

-- name: ListPackageTermsConditionsByPackageIdAsc :many
SELECT * FROM package_terms_conditions
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageTermsConditionsByPackageIdDesc :many
SELECT * FROM package_terms_conditions
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackageTermsConditions :one
INSERT INTO package_terms_conditions (
    id,
    created,
    modified,
    package_id,
    general_terms,
    luggage_notes,
    additional_notes
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
RETURNING *;

-- name: UpdatePackageTermsConditions :one
UPDATE package_terms_conditions
SET
    modified = now() AT TIME ZONE 'UTC',
    package_id = $2,
    general_terms = $3,
    luggage_notes = $4,
    additional_notes = $5
WHERE id = $1
RETURNING *;

-- name: DeletePackageTermsConditions :exec
DELETE FROM package_terms_conditions
WHERE id = $1;

-- name: DeletePackageTermsConditionsByPackageId :exec
DELETE FROM package_terms_conditions
WHERE package_id = $1;
