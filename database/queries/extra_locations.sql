-- name: GetExtraLocationById :one
SELECT * FROM extra_location
WHERE id = $1
LIMIT 1;

-- name: CountExtraLocationsByExtraId :one
SELECT
    COUNT(id)
FROM extra_location
WHERE extra_id = $1;

-- name: CountExtraLocationsByLocationId :one
SELECT
    COUNT(id)
FROM extra_location
WHERE location_id = $1;

-- name: ListExtraLocationsByExtraId :many
SELECT * FROM extra_location
WHERE extra_id = $1;

-- name: ListExtraLocationsByLocationId :many
SELECT * FROM extra_location
WHERE location_id = $1;

-- name: CreateExtraLocation :one
INSERT INTO extra_location (
    id,
    created,
    modified,
    extra_id,
    location_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: DeleteExtraLocation :exec
DELETE FROM extra_location
WHERE id = $1;

-- name: DeleteExtraLocationByExtraId :exec
DELETE FROM extra_location
WHERE extra_id = $1;

-- name: UpdateExtraLocation :one
UPDATE extra_location
SET
    modified = now() AT TIME ZONE 'UTC',
    extra_id = $2,
    location_id = $3
WHERE id = $1
RETURNING *;
