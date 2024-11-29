-- name: GetFareTypeById :one
SELECT * FROM fare_type
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetFareTypeByName :one
SELECT * FROM fare_type
WHERE name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountFareTypes :one
SELECT
    COUNT(id)
FROM fare_type
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: ListFareTypesAsc :many
SELECT * FROM fare_type
    WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListFareTypesDesc :many
SELECT * FROM fare_type
    WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateFareType :one
INSERT INTO fare_type (
    id,
    created,
    modified,
    deleted,
    name,
    ordering,
    description,
    is_international,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5
         )
RETURNING *;

-- name: DeleteFareType :exec
UPDATE fare_type
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateFareType :one
UPDATE fare_type
SET
    modified = now() AT TIME ZONE 'UTC',
    name = $2,
    ordering = $3,
    description = $4,
    is_international = $5,
    is_active = $6
WHERE id = $1
RETURNING *;

-- name: GetActiveFareTypes :many
SELECT * FROM fare_type
WHERE is_active = true
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY ordering;
