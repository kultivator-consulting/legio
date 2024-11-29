-- name: GetLocation :one
SELECT * FROM location
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetLocationById :one
SELECT * FROM location
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetLocationByName :one
SELECT * FROM location
WHERE name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: GetLocationByCode :one
SELECT * FROM location
WHERE code = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: ListStartLocationsByProductTypeSlug :many
SELECT DISTINCT location.*, station.latitude, station.longitude FROM service_route
  INNER JOIN product_type ON product_type.slug=$1
  INNER JOIN service ON service_route.service_id=service.id AND product_type.id=service.product_type_id
  INNER JOIN station ON service_route.start_station_id=station.id
  INNER JOIN location ON station.location_id=location.id
AND location.deleted > now() AT TIME ZONE 'UTC'
ORDER BY station.latitude, station.longitude;

-- name: ListEndLocationsByProductTypeSlug :many
SELECT DISTINCT location.*, station.latitude, station.longitude FROM service_route
  INNER JOIN product_type ON product_type.slug=$1
  INNER JOIN service ON service_route.service_id=service.id AND product_type.id=service.product_type_id
  INNER JOIN station ON service_route.end_station_id=station.id
  INNER JOIN location ON station.location_id=location.id
AND location.deleted > now() AT TIME ZONE 'UTC'
ORDER BY station.latitude, station.longitude;

-- name: ListPairedEndLocationsByStartLocationAndProductTypeSlug :many
SELECT DISTINCT location.*, end_station.latitude, end_station.longitude
FROM service_route
         INNER JOIN service ON service_route.service_id=service.id
         INNER JOIN product_type ON service.product_type_id=product_type.id
         INNER JOIN station start_station ON service_route.start_station_id=start_station.id
         INNER JOIN station end_station ON service_route.end_station_id=end_station.id
         INNER JOIN location ON location.id=end_station.location_id
WHERE product_type.slug=sqlc.arg(slug)::TEXT
  AND service.deleted > now() AT TIME ZONE 'UTC'
  AND start_station.location_id=sqlc.arg(start_location_id)::UUID
ORDER BY end_station.latitude, end_station.longitude;

-- name: ListPairedStartLocationsByEndLocationAndProductTypeSlug :many
SELECT DISTINCT location.*, start_station.latitude, start_station.longitude
FROM service_route
         INNER JOIN service ON service_route.service_id=service.id
         INNER JOIN product_type ON service.product_type_id=product_type.id
         INNER JOIN station start_station ON service_route.start_station_id=start_station.id
         INNER JOIN station end_station ON service_route.end_station_id=end_station.id
         INNER JOIN location ON location.id=start_station.location_id
WHERE product_type.slug=sqlc.arg(slug)::TEXT
  AND service.deleted > now() AT TIME ZONE 'UTC'
  AND end_station.location_id=sqlc.arg(end_location_id)::UUID
ORDER BY start_station.latitude, start_station.longitude;

-- name: CountLocations :one
SELECT
    COUNT(id)
FROM location
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListLocationsAsc :many
SELECT * FROM location
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListLocationsDesc :many
SELECT * FROM location
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateLocation :one
INSERT INTO location (
    id,
    created,
    modified,
    deleted,
    name,
    code,
    ordering,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
RETURNING *;

-- name: CreateLocationAndReturnId :one
INSERT INTO location (
    id,
    created,
    modified,
    deleted,
    name,
    code,
    ordering,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
RETURNING id;

-- name: DeleteLocation :exec
UPDATE location
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateLocation :one
UPDATE location
SET
    modified = now() AT TIME ZONE 'UTC',
    name = $2,
    code = $3,
    ordering = $4,
    is_active = $5
WHERE id = $1
RETURNING *;
