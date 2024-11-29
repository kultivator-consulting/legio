-- name: GetFareSearchResultById :one
SELECT * FROM fare_search_result
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: ListFareSearchResultsByFareSearchId :many
SELECT * FROM fare_search_result
WHERE fare_search_id = $1
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: CreateFareSearchResult :one
INSERT INTO fare_search_result (
    id,
    created,
    modified,
    deleted,
    fare_search_id,
    product_type_id,
    route_start_location,
    route_end_location,
    fare_date,
    is_return
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6
         )
RETURNING *;

-- name: DeleteFareSearchResult :exec
UPDATE fare_search_result
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateFareSearchResult :one
UPDATE fare_search_result
SET
    modified = now() AT TIME ZONE 'UTC',
    fare_search_id = $2,
    product_type_id = $3,
    route_start_location = $4,
    route_end_location = $5,
    fare_date = $6,
    is_return = $7
WHERE id = $1
RETURNING *;
