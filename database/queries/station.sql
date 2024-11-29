-- name: GetStationById :one
SELECT * FROM station
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetStationByName :one
SELECT * FROM station
WHERE name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: ListStationsByLocationId :many
SELECT * FROM station
WHERE location_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY name;

-- name: CountStations :one
SELECT
    COUNT(id)
FROM station
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListStationsAsc :many
SELECT * FROM station
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListStationsDesc :many
SELECT * FROM station
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateStation :one
INSERT INTO station (
    id,
    created,
    modified,
    deleted,
    location_id,
    name,
    address,
    code,
    description,
    latitude,
    longitude,
    stop_id,
    zone_id,
    stop_url,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
         )
RETURNING *;

-- name: DeleteStation :exec
UPDATE station
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateStation :one
UPDATE station
SET
    modified = now() AT TIME ZONE 'UTC',
    location_id = $2,
    name = $3,
    address = $4,
    code = $5,
    description = $6,
    latitude = $7,
    longitude = $8,
    stop_id = $9,
    zone_id = $10,
    stop_url = $11,
    is_active = $12
WHERE id = $1
RETURNING *;
