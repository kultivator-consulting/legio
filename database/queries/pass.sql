-- name: GetPassById :one
SELECT * FROM pass
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPassByCode :one
SELECT * FROM pass
WHERE code = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPassByName :one
SELECT * FROM pass
WHERE name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountPasses :one
SELECT
    COUNT(id)
FROM pass
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListPassesAsc :many
SELECT * FROM pass
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPassesDesc :many
SELECT * FROM pass
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CountPassesByParentPassId :one
SELECT
    COUNT(id)
FROM pass
WHERE parent_pass_id = $1
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: ListPassesByParentPassIdAsc :many
SELECT * FROM pass
WHERE parent_pass_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPassesByParentPassIdDesc :many
SELECT * FROM pass
WHERE parent_pass_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CountPassesAtRoot :one
SELECT
    COUNT(id)
FROM pass
WHERE parent_pass_id IS NULL
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: ListPassesAtRootAsc :many
SELECT * FROM pass
WHERE parent_pass_id IS NULL
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPassesAtRootDesc :many
SELECT * FROM pass
WHERE parent_pass_id IS NULL
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: GetChildPassesByParentPassId :many
SELECT * FROM pass
WHERE parent_pass_id = $1
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: CreatePass :one
INSERT INTO pass (
    id,
    created,
    modified,
    deleted,
    parent_pass_id,
    name,
    code,
    operator_codes,
    duration,
    duration_type,
    image,
    image_info,
    description,
    is_popular,
    is_top_up,
    ordering,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
         )
RETURNING *;

-- name: DeletePass :exec
UPDATE pass
SET deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdatePass :one
UPDATE pass
SET
    modified = now() AT TIME ZONE 'UTC',
    parent_pass_id = $2,
    name = $3,
    code = $4,
    operator_codes = $5,
    duration = $6,
    duration_type = $7,
    image = $8,
    image_info = $9,
    description = $10,
    is_popular = $11,
    is_top_up = $12,
    ordering = $13,
    is_active = $14
WHERE id = $1 AND deleted > now() AT TIME ZONE 'UTC'
RETURNING *;
