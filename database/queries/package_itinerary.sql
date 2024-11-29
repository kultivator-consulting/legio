-- name: GetPackageItineraryById :one
SELECT * FROM package_itinerary
WHERE id = $1
LIMIT 1;

-- name: GetPackageItinerariesByPackageId :many
SELECT * FROM package_itinerary
WHERE package_id = $1;

-- name: CountPackageItinerariesByPackageId :one
SELECT
    COUNT(id)
FROM package_itinerary
WHERE package_id = $1;

-- name: ListPackageItinerariesByPackageIdAsc :many
SELECT * FROM package_itinerary
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageItinerariesByPackageIdDesc :many
SELECT * FROM package_itinerary
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackageItinerary :one
INSERT INTO package_itinerary (
    id,
    created,
    modified,
    package_id,
    day,
    title,
    event_icon,
    station_id,
    latitude,
    longitude,
    description,
    supplier,
    supplier_code
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
         )
RETURNING *;

-- name: UpdatePackageItinerary :one
UPDATE package_itinerary
SET
    modified = now() AT TIME ZONE 'UTC',
    package_id = $2,
    day = $3,
    title = $4,
    event_icon = $5,
    station_id = $6,
    latitude = $7,
    longitude = $8,
    description = $9,
    supplier = $10,
    supplier_code = $11
WHERE id = $1
RETURNING *;

-- name: DeletePackageItinerary :exec
DELETE FROM package_itinerary
WHERE id = $1;

-- name: DeletePackageItineraryByPackageId :exec
DELETE FROM package_itinerary
WHERE package_id = $1;
