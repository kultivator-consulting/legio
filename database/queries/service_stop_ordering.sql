-- name: GetServiceStopOrderingById :one
SELECT * FROM service_stop_ordering
WHERE id = $1
LIMIT 1;

-- name: CountServiceStopOrderingsByServiceId :one
SELECT
    COUNT(id)
FROM service_stop_ordering
WHERE service_id = $1;

-- name: ListServiceStopOrderingsByServiceId :many
SELECT * FROM service_stop_ordering
WHERE service_id = $1
ORDER BY ordering;

-- name: CreateServiceStopOrdering :one
INSERT INTO service_stop_ordering (
    id,
    created,
    modified,
    service_id,
    location_id,
    ordering
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3
         )
RETURNING *;

-- name: DeleteServiceStopOrdering :exec
DELETE FROM service_stop_ordering
WHERE id = $1;

-- name: DeleteServiceStopOrderingByServiceId :exec
DELETE FROM service_stop_ordering
WHERE service_id = $1;

-- name: UpdateServiceStopOrdering :one
UPDATE service_stop_ordering
SET
    modified = now() AT TIME ZONE 'UTC',
    service_id = $2,
    location_id = $3,
    ordering = $4
WHERE id = $1
RETURNING *;
