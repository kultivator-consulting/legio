-- name: GetContentById :one
SELECT * FROM content
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetContentByDomainIdAndSlug :one
SELECT * FROM content
WHERE domain_id = $1
  AND slug = $2
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountContents :one
SELECT
    COUNT(id)
FROM content
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListContentsAsc :many
SELECT * FROM content
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListContentsDesc :many
SELECT * FROM content
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateContent :one
INSERT INTO content (
    id,
    created,
    modified,
    deleted,
    domain_id,
    component_id,
    account_id,
    title,
    slug,
    data,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7
         )
RETURNING *;

-- name: DeleteContent :exec
UPDATE content
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateContent :one
UPDATE content
SET
    modified = now() AT TIME ZONE 'UTC',
    domain_id = $2,
    component_id = $3,
    account_id = $4,
    title = $5,
    slug = $6,
    data = $7,
    is_active = $8
WHERE id = $1
RETURNING *;
