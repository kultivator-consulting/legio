-- name: GetPackageItineraryImageById :one
SELECT * FROM package_itinerary_images
WHERE id = $1
LIMIT 1;

-- name: GetPackageItineraryImagesByPackageItineraryId :many
SELECT * FROM package_itinerary_images
WHERE package_itinerary_id = $1;

-- name: CountPackageItineraryImagesByPackageItineraryId :one
SELECT
    COUNT(id)
FROM package_itinerary_images
WHERE package_itinerary_id = $1;

-- name: ListPackageItineraryImagesByPackageItineraryIdAsc :many
SELECT * FROM package_itinerary_images
WHERE package_itinerary_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageItineraryImagesByPackageItineraryIdDesc :many
SELECT * FROM package_itinerary_images
WHERE package_itinerary_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackageItineraryImage :one
INSERT INTO package_itinerary_images (
    id,
    created,
    modified,
    package_itinerary_id,
    ordering,
    image,
    image_info,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5
         )
RETURNING *;

-- name: UpdatePackageItineraryImage :one
UPDATE package_itinerary_images
SET
    modified = now() AT TIME ZONE 'UTC',
    package_itinerary_id = $2,
    ordering = $3,
    image = $4,
    image_info = $5,
    is_active = $6
WHERE id = $1
RETURNING *;

-- name: DeletePackageItineraryImage :exec
DELETE FROM package_itinerary_images
WHERE id = $1;

-- name: DeletePackageItineraryImagesByPackageItineraryId :exec
DELETE FROM package_itinerary_images
WHERE package_itinerary_id = $1;
