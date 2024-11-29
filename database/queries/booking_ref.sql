-- name: GetBookingRefById :one
SELECT * FROM booking_ref
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetBookingRefByPrefix :one
SELECT * FROM booking_ref
WHERE prefix = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountBookingRefs :one
SELECT
    COUNT(id)
FROM booking_ref
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListBookingRefsAsc :many
SELECT * FROM booking_ref
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListBookingRefsDesc :many
SELECT * FROM booking_ref
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateBookingRef :one
INSERT INTO booking_ref (
    id,
    created,
    modified,
    deleted,
    prefix,
    initial,
    length,
    sequence
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
RETURNING *;

-- name: DeleteBookingRef :exec
UPDATE booking_ref
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateBookingRef :one
UPDATE booking_ref
SET
    modified = now() AT TIME ZONE 'UTC',
    prefix = $2,
    initial = $3,
    length = $4,
    sequence = $5
WHERE id = $1
RETURNING *;

-- name: IncrementBookingRefSequenceById :one
UPDATE booking_ref
SET
    sequence = sequence + 1
WHERE id = $1
RETURNING *;

-- name: IncrementBookingRefSequenceByPrefix :one
UPDATE booking_ref
SET
    sequence = sequence + 1
WHERE prefix = $1
RETURNING *;
