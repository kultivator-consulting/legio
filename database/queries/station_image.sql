-- name: GetStationImageById :one
SELECT * FROM station_image
WHERE id = $1
LIMIT 1;

-- name: GetStationImagesByStationId :many
SELECT * FROM station_image
WHERE station_id = $1;

-- name: CountStationImagesByStationId :one
SELECT
    COUNT(id)
FROM station_image
WHERE station_id = $1;

-- name: ListStationImagesByStationIdAsc :many
SELECT * FROM station_image
WHERE station_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListStationImagesByStationIdDesc :many
SELECT * FROM station_image
WHERE station_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateStationImage :one
INSERT INTO station_image (
    id,
    created,
    modified,
    station_id,
    ordering,
    image,
    image_info
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
RETURNING *;

-- name: UpdateStationImage :one
UPDATE station_image
SET
    modified = now() AT TIME ZONE 'UTC',
    station_id = $2,
    ordering = $3,
    image = $4,
    image_info = $5
WHERE id = $1
RETURNING *;

-- name: DeleteStationImage :exec
DELETE FROM station_image
WHERE id = $1;

-- name: DeleteStationImagesByStationId :exec
DELETE FROM station_image
WHERE station_id = $1;
