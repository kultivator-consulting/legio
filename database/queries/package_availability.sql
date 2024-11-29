-- name: GetPackageAvailabilityById :one
SELECT * FROM package_availability
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPackageAvailabilityByPackageId :many
SELECT * FROM package_availability
WHERE package_id = $1;

-- name: GetValidPackageAvailabilityByDate :many
SELECT * FROM package_availability
WHERE package_id = $1
  AND start_date <= $2
  AND end_date >= $2;

-- name: CountPackageAvailabilityByPackageId :one
SELECT
    COUNT(id)
FROM package_availability
WHERE package_id = $1;

-- name: ListPackageAvailabilityByPackageIdAsc :many
SELECT * FROM package_availability
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageAvailabilityByPackageIdDesc :many
SELECT * FROM package_availability
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackageAvailability :one
INSERT INTO package_availability (
    id,
    created,
    modified,
    package_id,
    frequency,
    start_date,
    end_date,
    excluded_dates,
    sold_out_dates
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6
         )
RETURNING *;

-- name: UpdatePackageAvailability :one
UPDATE package_availability
SET
    modified = now() AT TIME ZONE 'UTC',
    package_id = $2,
    frequency = $3,
    start_date = $4,
    end_date = $5,
    excluded_dates = $6,
    sold_out_dates = $7
WHERE id = $1
RETURNING *;

-- name: DeletePackageAvailability :exec
DELETE FROM package_availability
WHERE id = $1;

-- name: DeletePackageAvailabilityByPackageId :exec
DELETE FROM package_availability
WHERE package_id = $1;
