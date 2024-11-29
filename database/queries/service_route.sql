-- name: GetServiceRouteById :one
SELECT * FROM service_route
WHERE id = $1
LIMIT 1;

-- name: GetServiceRouteByServiceId :one
SELECT * FROM service_route
WHERE service_id = $1
LIMIT 1;

-- name: CountServiceRoutesByServiceId :one
SELECT
    COUNT(id)
FROM service_route
WHERE service_id = $1;

-- name: ListServiceRoutesByServiceId :many
SELECT * FROM service_route
WHERE service_id = $1;

-- name: UpsertServiceRoute :one
INSERT INTO service_route (
    id,
    created,
    modified,
    service_id,
    product_id,
    start_station_id,
    end_station_id,
    is_main_route,
    is_popular,
    ic_route_id,
    short_name,
    description,
    url,
    color,
    text_color
) VALUES (
             coalesce(nullif(@id, uuid_nil()), uuid_generate_v4()),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
         )
ON CONFLICT (id) DO UPDATE
    SET
        modified = now() AT TIME ZONE 'UTC',
        service_id = $1,
        product_id = $2,
        start_station_id = $3,
        end_station_id = $4,
        is_main_route = $5,
        is_popular = $6,
        ic_route_id = $7,
        short_name = $8,
        description = $9,
        url = $10,
        color = $11,
        text_color = $12
RETURNING *;

-- name: CreateServiceRoute :one
INSERT INTO service_route (
    id,
    created,
    modified,
    service_id,
    product_id,
    start_station_id,
    end_station_id,
    is_main_route,
    is_popular,
    ic_route_id,
    short_name,
    description,
    url,
    color,
    text_color
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
         )
RETURNING *;

-- name: DeleteServiceRoute :exec
DELETE FROM service_route
WHERE id = $1;

-- name: DeleteServiceRouteByServiceId :exec
DELETE FROM service_route
WHERE service_id = $1;

-- name: UpdateServiceRoute :one
UPDATE service_route
SET
    modified = now() AT TIME ZONE 'UTC',
    service_id = $2,
    product_id = $3,
    start_station_id = $4,
    end_station_id = $5,
    is_main_route = $6,
    is_popular = $7,
    ic_route_id = $8,
    short_name = $9,
    description = $10,
    url = $11,
    color = $12,
    text_color = $13
WHERE id = $1
RETURNING *;
