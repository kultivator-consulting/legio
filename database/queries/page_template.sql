-- name: GetPageTemplateById :one
SELECT * FROM page_template
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPageTemplateAtRootAndBySlug :one
SELECT * FROM page_template
WHERE parent_page_path_id IS NULL
  AND slug = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPageTemplateByPagePathIdAndSlug :one
SELECT * FROM page_template
WHERE parent_page_path_id = $1
  AND slug = $2
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountPageTemplates :one
SELECT
    COUNT(id)
FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListAllPageTemplatesByPagePathIdAsc :many
SELECT * FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND parent_page_path_id = $1
ORDER BY title;

-- name: ListPageTemplatesByDomainIdAndAtRootAsc :many
SELECT * FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id IS NULL
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPageTemplatesByDomainIdAndAtRootDesc :many
SELECT * FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id IS NULL
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPageTemplatesByDomainIdAndPagePathIdAsc :many
SELECT * FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id = $2
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPageTemplatesByDomainIdAndPagePathIdDesc :many
SELECT * FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id = $2
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePageTemplateAndReturnId :one
INSERT INTO page_template (
    id,
    created,
    modified,
    deleted,
    domain_id,
    account_id,
    content_id,
    parent_page_path_id,
    title,
    slug,
    description,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8
         )
RETURNING id;

-- name: DeletePageTemplateById :exec
UPDATE page_template
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdatePageTemplateById :one
UPDATE page_template
SET
    modified = now() AT TIME ZONE 'UTC',
    domain_id = $2,
    account_id = $3,
    content_id = $4,
    parent_page_path_id = $5,
    title = $6,
    slug = $7,
    description = $8,
    is_active = $9
WHERE id = $1
RETURNING *;
