-- name: GetFareSearchResultAvailabilityById :one
SELECT * FROM fare_search_result_availability
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: ListFareSearchResultAvailabilityByFareSearchResultId :many
SELECT * FROM fare_search_result_availability
WHERE fare_search_result_id = $1
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: CreateFareSearchResultAvailability :one
INSERT INTO fare_search_result_availability (
    id,
    created,
    modified,
    deleted,
    fare_search_result_id,
    token,
    fare_date,
    available,
    availability_status,
    service_name,
    departure_city,
    departure_city_info,
    arrival_city,
    arrival_city_info,
    departure_time,
    arrival_time,
    fare_type,
    fare_type_description,
    adult_price,
    child_price,
    infant_price,
    description,
    ticket_conditions
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
         )
RETURNING *;

-- name: DeleteFareSearchResultAvailability :exec
UPDATE fare_search_result_availability
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateFareSearchResultAvailability :one
UPDATE fare_search_result_availability
SET
    modified = now() AT TIME ZONE 'UTC',
    fare_search_result_id = $2,
    token = $3,
    fare_date = $4,
    available = $5,
    availability_status = $6,
    service_name = $7,
    departure_city = $8,
    departure_city_info = $9,
    arrival_city = $10,
    arrival_city_info = $11,
    departure_time = $12,
    arrival_time = $13,
    fare_type = $14,
    fare_type_description = $15,
    adult_price = $16,
    child_price = $17,
    infant_price = $18,
    description = $19,
    ticket_conditions = $20
WHERE id = $1
RETURNING *;
