-- name: GetPassAvailabilityById :one
SELECT * FROM pass_availability
WHERE id = $1
LIMIT 1;

-- name: CountPassAvailabilities :one
SELECT
    COUNT(id)
FROM pass_availability
WHERE pass_id = $1;

-- name: ListPassAvailabilitiesByPassId :many
SELECT * FROM pass_availability
WHERE pass_id = $1;

-- name: CreatePassAvailability :one
INSERT INTO pass_availability (
    id,
    created,
    modified,
    pass_id,
    start_date,
    end_date,
    excluded_dates,
    notes,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6
         )
RETURNING *;

-- name: UpsertPassAvailability :one
INSERT INTO pass_availability (
    id,
    created,
    modified,
    pass_id,
    start_date,
    end_date,
    excluded_dates,
    notes,
    is_active
) VALUES (
             coalesce(nullif(@id, uuid_nil()), uuid_generate_v4()),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
                $1, $2, $3, $4, $5, $6
         )
ON CONFLICT (id) DO UPDATE
    SET
        modified = now() AT TIME ZONE 'UTC',
        pass_id = $1,
        start_date = $2,
        end_date = $3,
        excluded_dates = $4,
        notes = $5,
        is_active = $6
RETURNING *;

-- name: DeletePassAvailability :exec
DELETE FROM pass_availability
WHERE id = $1;

-- name: UpdatePassAvailability :one
UPDATE pass_availability
SET
    modified = now() AT TIME ZONE 'UTC',
    pass_id = $2,
    start_date = $3,
    end_date = $4,
    excluded_dates = $5,
    notes = $6,
    is_active = $7
WHERE id = $1
RETURNING *;
