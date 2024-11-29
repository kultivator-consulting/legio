-- name: GetFareSearchById :one
SELECT * FROM fare_search
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CreateFareSearch :one
INSERT INTO fare_search (
    id,
    created,
    modified,
    deleted,
    product_type_id,
    route_start_location_id,
    route_end_location_id,
    travel_date,
    is_return,
    adult_count,
    child_count,
    infant_count
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8
         )
RETURNING *;

-- name: DeleteFareSearch :exec
UPDATE fare_search
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateFareSearch :one
UPDATE fare_search
SET
    modified = now() AT TIME ZONE 'UTC',
    product_type_id = $2,
    route_start_location_id = $3,
    route_end_location_id = $4,
    travel_date = $5,
    is_return = $6,
    adult_count = $7,
    child_count = $8,
    infant_count = $9
WHERE id = $1
RETURNING *;
