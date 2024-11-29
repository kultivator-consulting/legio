-- name: GetServiceNoteById :one
SELECT * FROM service_note
WHERE id = $1
LIMIT 1;

-- name: GetServiceNoteByServiceId :one
SELECT * FROM service_note
WHERE service_id = $1
LIMIT 1;

-- name: CountServiceNotesByServiceId :one
SELECT
    COUNT(id)
FROM service_note
WHERE service_id = $1;

-- name: ListServiceNotesByServiceId :many
SELECT * FROM service_note
WHERE service_id = $1;

-- name: CreateServiceNote :one
INSERT INTO service_note (
    id,
    created,
    modified,
    service_id,
    note
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: DeleteServiceNote :exec
DELETE FROM service_note
WHERE id = $1;

-- name: DeleteServiceNoteByServiceId :exec
DELETE FROM service_note
WHERE service_id = $1;

-- name: UpdateServiceNote :one
UPDATE service_note
SET
    modified = now() AT TIME ZONE 'UTC',
    service_id = $2,
    note = $3
WHERE id = $1
RETURNING *;
