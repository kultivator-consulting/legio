-- name: GetExtraById :one
SELECT * FROM extra
WHERE id = $1 AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetExtraByName :one
SELECT * FROM extra
WHERE name = $1 AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetExtraByCode :one
SELECT * FROM extra
WHERE code = $1 AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountExtras :one
SELECT
    COUNT(id)
FROM extra
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListExtrasAsc :many
SELECT * FROM extra
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListExtrasDesc :many
SELECT * FROM extra
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateExtra :one
INSERT INTO extra (
    id,
    created,
    modified,
    deleted,
    name,
    banner_image,
    code,
    description,
    unit_price,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6
         )
RETURNING *;

-- name: DeleteExtra :exec
UPDATE extra
SET deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateExtra :one
UPDATE extra
SET
    modified = now() AT TIME ZONE 'UTC',
    name = $2,
    banner_image = $3,
    code = $4,
    description = $5,
    unit_price = $6,
    is_active = $7
WHERE id = $1 AND deleted > now() AT TIME ZONE 'UTC'
RETURNING *;
