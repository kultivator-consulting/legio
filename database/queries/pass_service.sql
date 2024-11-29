-- name: GetPassServiceById :one
SELECT * FROM pass_service
WHERE id = $1
LIMIT 1;

-- name: GetPassServiceByName :one
SELECT * FROM pass_service
WHERE name = $1
LIMIT 1;

-- name: CountPassServices :one
SELECT
    COUNT(id)
FROM pass_service
WHERE pass_id = $1;

-- name: ListPassServicesByPassId :many
SELECT * FROM pass_service
WHERE pass_id = $1;

-- name: CreatePassService :one
INSERT INTO pass_service (
    id,
    created,
    modified,
    pass_id,
    name,
    duration,
    duration_type
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
RETURNING *;

-- name: UpsertPassService :one
INSERT INTO pass_service (
    id,
    created,
    modified,
    pass_id,
    name,
    duration,
    duration_type
) VALUES (
             coalesce(nullif(@id, uuid_nil()), uuid_generate_v4()),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
ON CONFLICT (id) DO UPDATE
    SET
        modified = now() AT TIME ZONE 'UTC',
        pass_id = $1,
        name = $2,
        duration = $3,
        duration_type = $4
RETURNING *;

-- name: DeletePassService :exec
DELETE FROM pass_service
WHERE id = $1;

-- name: UpdatePassService :one
UPDATE pass_service
SET
    modified = now() AT TIME ZONE 'UTC',
    pass_id = $2,
    name = $3,
    duration = $4,
    duration_type = $5
WHERE id = $1
RETURNING *;
